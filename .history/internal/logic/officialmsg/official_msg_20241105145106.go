package officialmsg

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"gf-chat/internal/dao"
	"gf-chat/internal/model"
	"gf-chat/internal/model/entity"
	"gf-chat/internal/model/official"
	"gf-chat/internal/model/relation"
	"gf-chat/internal/service"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/gclient"
	"github.com/gogf/gf/v2/os/gcache"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/util/gconv"
)

func init() {
	service.RegisterOfficialMsg(&sOfficialMsg{})
}

type sOfficialMsg struct {
}

func (s sOfficialMsg) getOpenId(admin *relation.CustomerAdmins) string {
	wechat := entity.CustomerAdminWechat{}
	err := dao.CustomerAdminWechat.Ctx(gctx.New()).Where("admin_id", admin.Id).Scan(&wechat)
	if err == sql.ErrNoRows {
		return ""
	}
	return wechat.OfficialOpenId
}

func (s sOfficialMsg) getLimitKey(openId string) string {
	return fmt.Sprintf("official-message:%s", openId)
}

func (s sOfficialMsg) setLimit(ctx context.Context, openId string) {
	_ = gcache.Set(ctx, s.getLimitKey(openId), 1, time.Minute*5)
}

func (s sOfficialMsg) isLimit(ctx context.Context, openId string) bool {
	limit, _ := gcache.Get(ctx, s.getLimitKey(openId))
	if limit.Val() != nil {
		return true
	}
	return false
}

func (s sOfficialMsg) getUrl(ctx context.Context) string {
	url, _ := g.Cfg().Get(ctx, "app.ssoUrl")
	return url.String() + "/service/official-message"
}

func (s sOfficialMsg) Chat(admin *relation.CustomerAdmins, offline official.Offline) error {
	openId := s.getOpenId(admin)
	ctx := gctx.New()
	tmpl, err := g.Cfg().Get(ctx, "wechat.official.chatTmpl")
	if err != nil {
		return err
	}
	if openId == "" {
		return errors.New("用户没有绑定微信")
	}
	if s.isLimit(ctx, openId) {
		return errors.New("发送过于频繁")
	}
	message := official.Message{
		Miniprogram: nil,
		Data:        offline,
		Url:         "",
		TemplateId:  tmpl.String(),
		ToUser:      openId,
	}
	err = s.Send(gctx.New(), message)
	if err == nil {
		s.setLimit(ctx, openId)
		return nil
	}
	return err
}

func (s sOfficialMsg) Send(ctx context.Context, message official.Message) error {
	client := gclient.New()
	params := make(map[string]any)
	params["platform"] = "chat"
	params["message"] = gconv.Map(message)
	paramsByte, _ := json.Marshal(params)
	resp, err := client.Header(map[string]string{
		"Accept": "application/json",
	}).Post(ctx, s.getUrl(ctx), string(paramsByte))
	defer resp.Close()
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return errors.New(resp.Status)
	}
	body := resp.ReadAllString()
	r := model.ServiceResult{}
	err = json.Unmarshal([]byte(body), &r)
	if err != nil {
		return err
	}
	if r.Code == 0 {
		return nil
	}
	return errors.New(r.Message)
}

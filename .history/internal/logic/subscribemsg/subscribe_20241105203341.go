package subscribemsg

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"gf-chat/internal/dao"
	"gf-chat/internal/model"
	"gf-chat/internal/model/do"
	"gf-chat/internal/model/entity"
	"gf-chat/internal/service"
	"regexp"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/gclient"
	"github.com/gogf/gf/v2/os/gctx"
)

func init() {
	service.RegisterSubscribeMsg(&sSubscribeMsg{})
}

type sSubscribeMsg struct {
}

func (s sSubscribeMsg) getUrl(ctx context.Context) string {
	ssoUrl, _ := g.Cfg().Get(ctx, "app.saasUrl")
	return ssoUrl.String() + "/service/subscribe-message"
}

func (s sSubscribeMsg) Send(ctx context.Context, customerId, uid uint) error {
	if service.ChatSetting().GetSubscribeId(customerId) != "" {
		client := gclient.New()
		data := make(map[string]any)
		data["uid"] = uid
		data["platform"] = "chat"
		js, _ := json.Marshal(data)
		resp, err := client.Post(ctx, s.getUrl(ctx), string(js))
		defer resp.Close()
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
	return errors.New("subscribe tmpl not set")
}

func (s sSubscribeMsg) First(w do.WeappSubscribeMessages) *entity.WeappSubscribeMessages {
	item := &entity.WeappSubscribeMessages{}
	err := dao.WeappSubscribeMessages.Ctx(gctx.New()).Where(w).Scan(item)
	if err == sql.ErrNoRows {
		return nil
	}
	return item
}

func (s sSubscribeMsg) GetEntities(customerId uint) []entity.WeappSubscribeMessages {
	res := make([]entity.WeappSubscribeMessages, 0)
	_ = dao.WeappSubscribeMessages.Ctx(gctx.New()).
		Where("customer_id", customerId).Scan(&res)
	return res
}

func (s sSubscribeMsg) CheckChatTmpl(e entity.WeappSubscribeMessages) error {
	params := s.GetParams(e)
	length := len(params)
	if length != 2 {
		return errors.New("订阅消息模板必须只有2个关键字")
	}
	var thing, time bool
	for _, key := range params {
		if s.IsThing(key) {
			thing = true
		}
		if s.IsTime(key) {
			time = true
		}
	}
	if !thing || !time {
		return errors.New("关键字类型必须一个是time类型，一个是thing类型")
	}
	return nil
}

func (s sSubscribeMsg) IsTime(key string) bool {
	ok, _ := regexp.Match(`time\d`, []byte(key))
	return ok
}

func (s sSubscribeMsg) IsThing(key string) bool {
	ok, _ := regexp.Match(`thing\d`, []byte(key))
	return ok
}

func (s sSubscribeMsg) GetParams(e entity.WeappSubscribeMessages) []string {
	content := e.Content
	reg := regexp.MustCompile(`\{\{(.*)\.DATA}}`)
	matchs := reg.FindAllStringSubmatch(content, -1)
	params := make([]string, 0)
	for _, item := range matchs {
		if len(item) >= 2 {
			params = append(params, item[1])
		}
	}
	return params
}

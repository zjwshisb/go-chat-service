package sms

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
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/gclient"
	"github.com/gogf/gf/v2/os/gctx"
	"regexp"
)

func init() {
	service.RegisterSms(&sSms{})
}

type sSms struct {
}

func (sSms sSms) getUrl(ctx context.Context) string {
	ssoUrl, _ := g.Cfg().Get(ctx, "app.saasUrl")
	return ssoUrl.String() + "/service/sms-message"
}

func (sSms sSms) Send(ctx context.Context, code string, uid int) error {
	client := gclient.New()
	params := make(map[string]any)
	params["platform"] = "chat"
	params["code"] = code
	params["uid"] = uid
	paramsByte, _ := json.Marshal(params)
	resp, err := client.Header(map[string]string{
		"Accept": "application/json",
	}).Post(ctx, sSms.getUrl(ctx), string(paramsByte))
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

func (sSms sSms) First(w do.SmsTemplates) *entity.SmsTemplates {
	model := &entity.SmsTemplates{}
	err := dao.SmsTemplates.Ctx(gctx.New()).Where(w).Scan(model)
	if err == sql.ErrNoRows {
		return nil
	}
	return model
}

func (sSms *sSms) CheckChatSms(model *entity.SmsTemplates) error {
	content := model.Content
	hasParam, _ := regexp.Match(`\$\{(.*)}`, []byte(content))
	if hasParam {
		return errors.New("不支持带有变量的模板")
	}
	return nil
}

func (sSms *sSms) GetValidTemplate(customerId int) []entity.SmsTemplates {
	res := make([]entity.SmsTemplates, 0)
	dao.SmsTemplates.Ctx(gctx.New()).Where(do.SmsTemplates{
		CustomerId: customerId,
		Status:     1,
	}).Scan(&res)
	return res
}

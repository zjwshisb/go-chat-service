package sso

import (
	"context"
	"encoding/json"
	"errors"
	"gf-chat/internal/model"
	"gf-chat/internal/service"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/util/gconv"
	"net/http"
)

func init() {
	service.RegisterSso(New())
}

func New() *sSso {
	return &sSso{}
}

type sSso struct {
}

type result struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		Id        int    `json:"id"`
		SessionId string `json:"session_id"`
	} `json:"data"`
}

func getUrl() string {
	ctx := gctx.New()
	ssoUrl, _ := g.Cfg().Get(ctx, "app.ssoUrl")
	return ssoUrl.String()
}

func (sSso *sSso) Check(ctx context.Context, sessionId string, uid int) bool {
	data := make(map[string]string)
	data["session_id"] = sessionId
	js, _ := json.Marshal(data)
	client := g.Client()
	resp, err := client.Post(ctx, getUrl()+"/sso/user", string(js))
	if err != nil {
		return false
	}
	defer resp.Close()
	body := resp.ReadAllString()
	result := model.ServiceResult{}
	err = json.Unmarshal([]byte(body), &result)
	if err != nil {
		return false
	}
	return result.Code == 0 && gconv.Int(result.Data) == uid
}

func (sSso *sSso) Auth(ctx context.Context, ticket string) (uid int, sessionId string, err error) {
	data := make(map[string]string)
	data["ticket"] = ticket
	client := g.Client()
	resp, err := client.Post(ctx, getUrl()+"/sso/auth", data)
	body := resp.ReadAllString()
	if resp.StatusCode != http.StatusOK {
		err = errors.New(resp.Status)
		return
	} else {
		if err != nil {
			return
		}
		result := &result{}
		err = json.Unmarshal([]byte(body), result)
		if err != nil {
			return
		}
		if result.Message != "" {
			err = errors.New(result.Message)
			return
		}
		return result.Data.Id, result.Data.SessionId, nil
	}

}

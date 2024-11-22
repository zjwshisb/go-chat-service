package automessage

import (
	"encoding/json"
	api "gf-chat/api/v1/backend"
	"gf-chat/internal/consts"
	"gf-chat/internal/dao"
	"gf-chat/internal/model"
	"gf-chat/internal/model/entity"
	"gf-chat/internal/service"
	"gf-chat/internal/trait"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
)

func init() {
	service.RegisterAutoMessage(&sAutoMessage{
		Curd: trait.Curd[model.CustomerChatAutoMessage]{
			Dao: &dao.CustomerChatAutoMessages,
		},
	})
}

type sAutoMessage struct {
	trait.Curd[model.CustomerChatAutoMessage]
}

func (s *sAutoMessage) Fill(message *model.CustomerChatAutoMessage, form api.AutoMessageForm) {
	message.Name = form.Name
	message.Type = form.Type
	switch message.Type {
	case consts.MessageTypeNavigate:
		content := g.Map{
			"title": form.Navigator.Title,
			"image": form.Navigator.Image.Id,
			"url":   form.Navigator.Url,
		}
		contentJson, _ := json.Marshal(content)
		message.Content = string(contentJson)
	case consts.MessageTypeText:
		message.Content = form.Content
	case consts.MessageTypeFile:
		message.Content = gconv.String(form.File.Id)
	default:
	}
}
func (s *sAutoMessage) ToChatMessage(auto *model.CustomerChatAutoMessage) (msg *model.CustomerChatMessage, err error) {
	content := auto.Content
	if auto.Type == consts.MessageTypeFile {
		//content = service.Qiniu().Url(content)
	}
	if auto.Type == consts.MessageTypeNavigate {
		m := make(map[string]string)
		err = json.Unmarshal([]byte(auto.Content), &m)
		if err != nil {
			return
		}
		//m["content"] = service.Qiniu().Url(m["content"])
		newT, err := json.Marshal(m)
		if err != nil {
			return nil, err
		}
		content = string(newT)
	}

	return &model.CustomerChatMessage{
		CustomerChatMessages: entity.CustomerChatMessages{
			UserId:     0,
			AdminId:    0,
			Type:       auto.Type,
			Content:    content,
			CustomerId: auto.CustomerId,
			Source:     consts.MessageSourceSystem,
			SessionId:  0,
			ReqId:      service.ChatMessage().GenReqId(),
		},
	}, err
}

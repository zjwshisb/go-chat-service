package automessage

import (
	"context"
	"encoding/json"
	"gf-chat/api/v1/backend/automessage"
	"gf-chat/internal/consts"
	"gf-chat/internal/dao"
	"gf-chat/internal/model"
	"gf-chat/internal/model/entity"
	"gf-chat/internal/service"
	"gf-chat/internal/trait"
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

func (s *sAutoMessage) UpdateOne(ctx context.Context, message *model.CustomerChatAutoMessage, req *automessage.UpdateReq) (count int64, err error) {
	message.Name = req.Name
	switch message.Type {
	case consts.MessageTypeNavigate:
		content := map[string]string{
			"title":   req.Title,
			"url":     req.Url,
			"content": req.Content,
		}
		contentJson, _ := json.Marshal(content)
		message.Content = string(contentJson)
	default:
		message.Content = req.Content
	}
	result, err := dao.CustomerChatAutoMessages.Ctx(ctx).Save(message)
	if err != nil {
		return
	}
	return result.RowsAffected()
}

func (s *sAutoMessage) SaveOne(ctx context.Context, req *automessage.StoreReq) (id int64, err error) {
	admin := service.AdminCtx().GetAdmin(ctx)
	item := entity.CustomerChatAutoMessages{
		Name:       req.Name,
		Type:       req.Type,
		CustomerId: admin.CustomerId,
	}
	switch item.Type {
	case consts.MessageTypeNavigate:
		content := map[string]string{
			"title":   req.Title,
			"url":     req.Url,
			"content": req.Content,
		}
		contentJson, _ := json.Marshal(content)
		item.Content = string(contentJson)
	default:
		item.Content = req.Content
	}
	result, err := dao.CustomerChatAutoMessages.Ctx(ctx).Insert(&item)
	if err != nil {
		return
	}
	return result.LastInsertId()
}

func (s *sAutoMessage) ToChatMessage(auto *entity.CustomerChatAutoMessages) (msg *model.CustomerChatMessage, err error) {
	content := auto.Content
	if auto.Type == consts.MessageTypeImage {
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

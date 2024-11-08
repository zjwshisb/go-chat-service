package automessage

import (
	"context"
	"database/sql"
	"encoding/json"
	"gf-chat/api/v1/backend/automessage"
	"gf-chat/internal/consts"
	"gf-chat/internal/dao"
	"gf-chat/internal/model"
	"gf-chat/internal/model/do"
	"gf-chat/internal/model/entity"
	"gf-chat/internal/service"
)

func init() {
	service.RegisterAutoMessage(&sAutoMessage{})
}

type sAutoMessage struct {
}

func (s *sAutoMessage) First(ctx context.Context, w any) (msg *entity.CustomerChatAutoMessages, err error) {
	err = dao.CustomerChatAutoMessages.Ctx(ctx).Where(w).Scan(&msg)
	if err != nil {
		return
	}
	if msg == nil {
		err = sql.ErrNoRows
	}
	return
}
func (s *sAutoMessage) Paginate(ctx context.Context, w *do.CustomerChatAutoMessages, p model.QueryInput) (items []*entity.CustomerChatAutoMessages, total int) {
	query := dao.CustomerChatAutoMessages.Ctx(ctx)
	if w != nil {
		query = query.Where(w)
	}
	// if p.WithTotal {
	// 	total, _ = query.Clone().Count()
	// 	if total == 0 {
	// 		return
	// 	}
	// }
	err := query.Page(p.GetPage(), p.GetSize()).Scan(&items)
	if err == sql.ErrNoRows {
		return
	}
	if total == 0 {
		total = len(items)
	}
	return
}

func (s *sAutoMessage) All(ctx context.Context, w do.CustomerChatAutoMessages) (items []*entity.CustomerChatAutoMessages, err error) {
	err = dao.CustomerChatAutoMessages.Ctx(ctx).Where(w).Scan(&items)
	if items == nil {
		err = sql.ErrNoRows
	}
	if err != nil {
		items = make([]*entity.CustomerChatAutoMessages, 0)
	}
	return

}

func (s *sAutoMessage) EntityToListItem(i entity.CustomerChatAutoMessages) model.AutoMessageListItem {
	l := model.AutoMessageListItem{
		Id:         i.Id,
		Name:       i.Name,
		Type:       i.Type,
		Content:    i.Content,
		CreatedAt:  i.CreatedAt,
		UpdatedAt:  i.UpdatedAt,
		RulesCount: 0,
	}
	if i.Type == consts.MessageTypeImage {
		l.Content = service.Qiniu().Url(i.Content)
	}
	if i.Type == consts.MessageTypeNavigate {
		m := make(map[string]string)
		_ = json.Unmarshal([]byte(i.Content), &m)
		l.Title = m["title"]
		l.Content = service.Qiniu().Url(m["content"])
		l.Url = m["url"]
	}
	return l
}

func (s *sAutoMessage) Update(ctx context.Context, message *entity.CustomerChatAutoMessages, req *automessage.UpdateReq) (count int64, err error) {
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

func (s *sAutoMessage) Save(ctx context.Context, req *automessage.StoreReq) (id int64, err error) {
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
		content = service.Qiniu().Url(content)
	}
	if auto.Type == consts.MessageTypeNavigate {
		m := make(map[string]string)
		err = json.Unmarshal([]byte(auto.Content), &m)
		if err != nil {
			return
		}
		m["content"] = service.Qiniu().Url(m["content"])
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

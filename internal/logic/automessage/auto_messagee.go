package automessage

import (
	"context"
	"encoding/json"
	api "gf-chat/api/v1/backend"
	"gf-chat/internal/consts"
	"gf-chat/internal/dao"
	"gf-chat/internal/model"
	"gf-chat/internal/model/do"
	"gf-chat/internal/model/entity"
	"gf-chat/internal/service"
	"gf-chat/internal/trait"
	"github.com/duke-git/lancet/v2/maputil"
	"github.com/duke-git/lancet/v2/slice"
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

func (s *sAutoMessage) ToApi(ctx context.Context, message *model.CustomerChatAutoMessage,
	files *map[uint]*model.CustomerChatFile) *api.AutoMessage {
	resp := &api.AutoMessage{
		Id:        message.Id,
		Name:      message.Name,
		Type:      message.Type,
		Content:   message.Content,
		CreatedAt: message.CreatedAt,
		UpdatedAt: message.UpdatedAt,
	}
	if message.Type == consts.MessageTypeFile {
		if files != nil {
			file := maputil.GetOrDefault(*files, gconv.Uint(message.Content), nil)
			if file != nil {
				resp.File = service.File().ToApi(file)
			}
		} else {
			resp.File, _ = service.File().FindAnd2Api(ctx, message.Content)
		}
	}
	if message.Type == consts.MessageTypeNavigate {
		var simpleNavigator *api.AutoMessageOriginNavigator
		err := json.Unmarshal([]byte(message.Content), &simpleNavigator)
		if err == nil {
			navigator := api.AutoMessageNavigator{}
			navigator.Url = simpleNavigator.Url
			navigator.Title = simpleNavigator.Title
			if files != nil {
				image := maputil.GetOrDefault(*files, simpleNavigator.Image, nil)
				if image != nil {
					navigator.Image = service.File().ToApi(image)
				}
			} else {
				navigator.Image, _ = service.File().FindAnd2Api(ctx, navigator.Image)
			}
			resp.Navigator = &navigator
		}
	}
	return resp
}

func (s *sAutoMessage) ToApis(ctx context.Context, items []*model.CustomerChatAutoMessage) (resp []*api.AutoMessage, err error) {
	resp = make([]*api.AutoMessage, len(items))
	filesId := slice.Map(items, func(index int, item *model.CustomerChatAutoMessage) any {
		switch item.Type {
		case consts.MessageTypeFile:
			return item.Content
		case consts.MessageTypeNavigate:
			var navigator *api.AutoMessageOriginNavigator
			err := json.Unmarshal([]byte(item.Content), &navigator)
			if err != nil {
				return 0
			}
			return navigator.Image
		}
		return 0
	})
	files, err := service.File().All(ctx, do.CustomerChatFiles{
		Id: slice.Unique(filesId),
	}, nil, nil)
	filesMap := slice.KeyBy(files, func(item *model.CustomerChatFile) uint {
		return item.Id
	})
	if err != nil {
		return
	}
	for index, i := range items {
		item := s.ToApi(ctx, i, &filesMap)
		resp[index] = item
	}
	return
}

func (s *sAutoMessage) Form2Do(form api.AutoMessageForm) *do.CustomerChatAutoMessages {
	message := &do.CustomerChatAutoMessages{}
	message.Name = form.Name
	message.Type = form.Type
	// UpdatedAt 不为nil时不会自动更新时间
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
	return message
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

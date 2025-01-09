package automessage

import (
	"context"
	"encoding/json"
	api "gf-chat/api/backend/v1"
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
	files *map[uint]*model.CustomerChatFile) (*api.AutoMessage, error) {
	resp := &api.AutoMessage{
		Id:        message.Id,
		Name:      message.Name,
		Type:      message.Type,
		Content:   message.Content,
		CreatedAt: message.CreatedAt,
		UpdatedAt: message.UpdatedAt,
	}
	if service.ChatMessage().IsFileType(message.Type) {
		if files != nil {
			file := maputil.GetOrDefault(*files, gconv.Uint(message.Content), nil)
			if file != nil {
				resp.File = service.File().ToApi(file)
			}
		} else {
			var err error
			resp.File, err = service.File().FindAnd2Api(ctx, message.Content)
			if err != nil {
				return nil, err
			}
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
				navigator.Image, _ = service.File().FindAnd2Api(ctx, simpleNavigator.Image)
			}
			resp.Navigator = &navigator
		}
	}
	return resp, nil
}

func (s *sAutoMessage) ToApis(ctx context.Context, items []*model.CustomerChatAutoMessage) (resp []*api.AutoMessage, err error) {
	resp = make([]*api.AutoMessage, len(items))
	filesId := slice.Map(items, func(index int, item *model.CustomerChatAutoMessage) any {
		if service.ChatMessage().IsFileType(item.Type) {
			return item.Content
		}
		if item.Type == consts.MessageTypeNavigate {
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
		item, err := s.ToApi(ctx, i, &filesMap)
		if err != nil {
			return resp, err
		}
		resp[index] = item
	}
	return
}

func (s *sAutoMessage) Form2Do(form api.AutoMessageForm) *do.CustomerChatAutoMessages {
	message := &do.CustomerChatAutoMessages{}
	message.Name = form.Name
	message.Type = form.Type
	message.Content = ""
	if service.ChatMessage().IsFileType(form.Type) {
		message.Content = gconv.String(form.File.Id)
	} else if form.Type == consts.MessageTypeNavigate {
		content := g.Map{
			"title": form.Navigator.Title,
			"image": form.Navigator.Image.Id,
			"url":   form.Navigator.Url,
		}
		contentJson, _ := json.Marshal(content)
		message.Content = string(contentJson)
	} else if form.Type == consts.ChatSettingTypeText {
		message.Content = form.Content
	}

	return message
}
func (s *sAutoMessage) ToChatMessage(ctx context.Context, auto *model.CustomerChatAutoMessage) (msg *model.CustomerChatMessage, err error) {
	apiMessage, err := s.ToApi(ctx, auto, nil)
	if err != nil {
		return
	}
	chatMessage := &model.CustomerChatMessage{
		CustomerChatMessages: entity.CustomerChatMessages{
			UserId:     0,
			AdminId:    0,
			Type:       apiMessage.Type,
			CustomerId: auto.CustomerId,
			Source:     consts.MessageSourceSystem,
			SessionId:  0,
			ReqId:      service.ChatMessage().GenReqId(),
		},
	}
	if service.ChatMessage().IsFileType(auto.Type) && apiMessage.File != nil {
		chatMessage.Content = apiMessage.File.Url
	} else if auto.Type == consts.MessageTypeNavigate {
		if apiMessage.Navigator != nil {
			mapContent := g.Map{
				"title": apiMessage.Navigator.Title,
				"url":   apiMessage.Navigator.Url,
			}
			if apiMessage.Navigator.Image != nil {
				mapContent["image"] = apiMessage.Navigator.Image.Url
			}
			content, _ := json.Marshal(mapContent)
			chatMessage.Content = string(content)
		}
	} else if auto.Type == consts.ChatSettingTypeText {
		chatMessage.Content = apiMessage.Content
	}

	return chatMessage, nil
}

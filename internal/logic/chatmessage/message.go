package chatmessage

import (
	"context"
	api "gf-chat/api/v1/backend"
	"gf-chat/internal/consts"
	"gf-chat/internal/dao"
	"gf-chat/internal/model"
	"gf-chat/internal/model/do"
	"gf-chat/internal/model/entity"
	"gf-chat/internal/service"
	"gf-chat/internal/trait"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/gogf/gf/v2/util/grand"
)

func init() {
	service.RegisterChatMessage(&sChatMessage{
		Curd: trait.Curd[model.CustomerChatMessage]{
			Dao: &dao.CustomerChatMessages,
		},
	})
}

type sChatMessage struct {
	trait.Curd[model.CustomerChatMessage]
}

func (s *sChatMessage) GenReqId() string {
	return grand.S(20)
}

func (s *sChatMessage) SaveWithUpdate(ctx context.Context, msg *model.CustomerChatMessage) (err error) {
	id, err := s.Save(ctx, msg)
	if err != nil {
		return
	}
	if msg.Id == 0 {
		msg.Id = uint(id)
	}
	return
}

func (s *sChatMessage) ToRead(ctx context.Context, id any) (int64, error) {
	return s.Update(ctx, g.Map{
		"read_at": nil,
		"id":      id,
	}, do.CustomerChatMessages{
		ReadAt: gtime.Now(),
	})
}

func (s *sChatMessage) GetAdminName(ctx context.Context, model *model.CustomerChatMessage) (avatar string, err error) {
	switch model.Source {
	case consts.MessageSourceAdmin:
		if model.Admin.Setting != nil && model.Admin.Setting.Name != "" {
			return model.Admin.Setting.Name, nil
		}
		avatar, err = service.ChatSetting().GetName(ctx, model.CustomerId)
		return
	case consts.MessageSourceSystem:
		avatar, err = service.ChatSetting().GetName(ctx, model.CustomerId)
		return
	}
	return "", nil
}
func (s *sChatMessage) ToApi(ctx context.Context, message *model.CustomerChatMessage) (msg *api.ChatMessage, err error) {
	username := ""
	if message.User != nil {
		username = message.User.Username
	}
	avatar, err := s.GetAvatar(ctx, message)
	if err != nil {
		return
	}
	name, err := s.GetAdminName(ctx, message)
	if err != nil {
		return
	}
	msg = &api.ChatMessage{
		Id:         message.Id,
		UserId:     message.UserId,
		AdminId:    message.AdminId,
		AdminName:  name,
		Type:       message.Type,
		Content:    message.Content,
		ReceivedAT: message.ReceivedAt,
		Source:     message.Source,
		ReqId:      message.ReqId,
		IsSuccess:  true,
		IsRead:     message.ReadAt != nil,
		Avatar:     avatar,
		Username:   username,
	}
	return
}
func (s *sChatMessage) GetAvatar(ctx context.Context, model *model.CustomerChatMessage) (avatar string, err error) {
	switch model.Source {
	case consts.MessageSourceAdmin:
		if model.Admin != nil &&
			model.Admin.Setting != nil {
			return service.Admin().GetAvatar(ctx, model.Admin)
		} else {
			return service.ChatSetting().GetAvatar(ctx, model.CustomerId)
		}
	case consts.MessageSourceSystem:
		return service.ChatSetting().GetAvatar(ctx, model.CustomerId)
	case consts.MessageSourceUser:
		return "", nil
	}
	return "", nil
}

func (s *sChatMessage) GetList(ctx context.Context, lastId uint, w any, size uint) (res []*model.CustomerChatMessage, err error) {
	res, err = service.ChatMessage().All(ctx, do.CustomerChatMessages{
		UserId: service.UserCtx().GetId(ctx),
	}, g.Slice{model.CustomerChatMessage{}.Admin,
		model.CustomerChatMessage{}.User}, "id desc", 20)
	//if size > 0 {
	//	query = query.Limit(int(size))
	//}
	//if lastId > 0 {
	//	query = query.Where("id < ?", lastId)
	//}
	//err = query.Scan(&res)
	//if err != nil {
	//	if errors.Is(err, sql.ErrNoRows) {
	//		res = make([]*model.CustomerChatMessage, 0)
	//	} else {
	//		return
	//	}
	//}
	return
}

func (s *sChatMessage) NewNotice(session *model.CustomerChatSession, content string) *model.CustomerChatMessage {
	return &model.CustomerChatMessage{
		CustomerChatMessages: entity.CustomerChatMessages{
			UserId:     session.UserId,
			AdminId:    session.AdminId,
			Type:       consts.MessageTypeNotice,
			Content:    content,
			CustomerId: session.CustomerId,
			SessionId:  session.Id,
			ReqId:      s.GenReqId(),
			Source:     consts.MessageSourceSystem,
		},
	}
}
func (s *sChatMessage) NewOffline(admin *model.CustomerAdmin) *model.CustomerChatMessage {
	if admin.Setting != nil && admin.Setting.OfflineContent != "" {
		return &model.CustomerChatMessage{
			CustomerChatMessages: entity.CustomerChatMessages{
				UserId:     0,
				AdminId:    admin.Id,
				Type:       consts.MessageTypeText,
				Content:    admin.Setting.OfflineContent,
				ReceivedAt: gtime.New(),
				CustomerId: admin.CustomerId,
				Source:     consts.MessageSourceAdmin,
				ReqId:      service.ChatMessage().GenReqId(),
			},
			Admin: admin,
			User:  nil,
		}
	}
	return nil
}

func (s *sChatMessage) NewWelcome(admin *model.CustomerAdmin) *model.CustomerChatMessage {
	if admin.Setting != nil && admin.Setting.WelcomeContent != "" {
		return &model.CustomerChatMessage{
			CustomerChatMessages: entity.CustomerChatMessages{
				UserId:     0,
				AdminId:    admin.Id,
				Type:       consts.MessageTypeText,
				Content:    admin.Setting.WelcomeContent,
				ReceivedAt: gtime.New(),
				CustomerId: admin.CustomerId,
				Source:     consts.MessageSourceAdmin,
				ReqId:      service.ChatMessage().GenReqId(),
			},
			Admin: admin,
			User:  nil,
		}
	}
	return nil
}

func (s *sChatMessage) Insert(ctx context.Context, message *model.CustomerChatMessage) (*model.CustomerChatMessage, error) {
	id, err := s.Save(ctx, message)
	if err != nil {
		return nil, err
	}
	message.Id = gconv.Uint(id)
	return message, nil
}

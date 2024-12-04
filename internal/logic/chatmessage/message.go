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
	"github.com/duke-git/lancet/v2/slice"
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

func (s *sChatMessage) GetUnreadCountGroupByUsers(ctx context.Context, uids []uint, w any) (res []model.UnreadCount, err error) {
	res = make([]model.UnreadCount, 0)
	err = s.Dao.Ctx(ctx).Where("user_id in (?)", uids).
		FieldCount("*", "Count").
		Fields("user_id").
		Group("user_id").
		Where(w).
		Where("read_at", nil).
		Scan(&res)
	return
}

func (s *sChatMessage) GetLastGroupByUsers(ctx context.Context, adminId uint, uids []uint) (res []*model.CustomerChatMessage, err error) {
	type idStruct struct {
		Id uint
	}
	var idArr []idStruct
	err = s.Dao.Ctx(ctx).
		Fields("max(id) as id").
		Where("user_id in (?)", uids).
		Where("source in (?)", []int{consts.MessageSourceAdmin, consts.MessageSourceUser}).
		Where("admin_id", adminId).
		Group("user_id").Scan(&idArr)
	if err != nil {
		res = make([]*model.CustomerChatMessage, 0)
		return
	}
	res, err = s.All(ctx, do.CustomerChatMessages{
		Id: slice.Map(idArr, func(_ int, i idStruct) uint {
			return i.Id
		}),
	}, nil, nil)
	return
}

func (s *sChatMessage) GenReqId() string {
	return grand.S(20)
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
		if model.Admin != nil && model.Admin.Setting != nil && model.Admin.Setting.Name != "" {
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
				ReceivedAt: gtime.Now(),
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
				ReceivedAt: gtime.Now(),
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

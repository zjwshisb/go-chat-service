package chatmessage

import (
	"context"
	"database/sql"
	"gf-chat/internal/consts"
	"gf-chat/internal/dao"
	"gf-chat/internal/model"
	"gf-chat/internal/model/entity"
	"gf-chat/internal/service"
	"gf-chat/internal/trait"
	"time"

	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/gogf/gf/v2/util/grand"
)

func init() {
	service.RegisterChatMessage(&sChatMessage{
		Curd: trait.Curd[model.CustomerChatMessage]{
			Dao: dao.CustomerChatMessages,
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
func (s *sChatMessage) EntityToRelation(msg *entity.CustomerChatMessages) *model.CustomerChatMessage {
	return &model.CustomerChatMessage{
		CustomerChatMessages: *msg,
		Admin:                nil,
		User:                 nil,
	}
}

func (s *sChatMessage) ChangeToRead(msgId []uint) (sql.Result, error) {
	ctx := gctx.New()
	query := dao.CustomerChatMessages.Ctx(ctx).Data(map[string]int64{
		"read_at": time.Now().Unix(),
	}).WhereIn("id", msgId).Where("read_at", 0)
	return query.Update()
}

func (s *sChatMessage) GetAdminName(model model.CustomerChatMessage) string {
	switch model.Source {
	case consts.MessageSourceAdmin:
		if model.Admin.Setting != nil && model.Admin.Setting.Name != "" {
			return model.Admin.Setting.Name
		}
		return service.ChatSetting().GetName(model.CustomerId)
	case consts.MessageSourceSystem:
		return service.ChatSetting().GetName(model.CustomerId)
	}
	return ""
}
func (s *sChatMessage) RelationToChat(message model.CustomerChatMessage) model.ChatMessage {
	username := ""
	if message.User != nil {
		username = message.User.Username
	}
	return model.ChatMessage{
		Id:         message.Id,
		UserId:     message.UserId,
		AdminId:    message.AdminId,
		AdminName:  s.GetAdminName(message),
		Type:       message.Type,
		Content:    message.Content,
		ReceivedAT: message.ReceivedAt,
		Source:     message.Source,
		ReqId:      message.ReqId,
		IsSuccess:  true,
		IsRead:     message.ReadAt != nil,
		Avatar:     s.GetAvatar(message),
		Username:   username,
	}
}
func (s *sChatMessage) GetAvatar(model model.CustomerChatMessage) string {
	switch model.Source {
	case consts.MessageSourceAdmin:
		if model.Admin != nil &&
			model.Admin.Setting != nil &&
			model.Admin.Setting.Avatar != "" {
			return service.Qiniu().Url(model.Admin.Setting.Avatar)
		} else {
			return service.ChatSetting().GetAvatar(model.CustomerId)
		}
	case consts.MessageSourceSystem:
		return service.ChatSetting().GetAvatar(model.CustomerId)
	case consts.MessageSourceUser:
		return ""
	}
	return ""
}

func (s *sChatMessage) GetModels(lastId uint, w any, size uint) []*model.CustomerChatMessage {
	res := make([]*model.CustomerChatMessage, 0)
	ctx := gctx.New()
	query := dao.CustomerChatMessages.Ctx(ctx).With(
		model.CustomerChatMessage{}.Admin,
		// todo
		// 多层关联会用N+1问题
		// relation.CustomerAdmins{}.Setting,
		model.CustomerChatMessage{}.User,
	).Where(w).OrderDesc("id")
	if size > 0 {
		query = query.Limit(int(size))
	}
	if lastId > 0 {
		query = query.Where("id < ?", lastId)
	}
	query.Scan(&res)
	adminMaps := make(map[uint]*model.CustomerAdmin)
	for _, message := range res {
		if message.Admin != nil {
			if _, ok := adminMaps[message.Admin.Id]; !ok {
				adminMaps[message.Admin.Id] = message.Admin
			}
		}
	}
	adminIds := make([]uint, 0, len(adminMaps))
	for _, admin := range adminMaps {
		adminIds = append(adminIds, admin.Id)
	}
	settings := make([]*entity.CustomerAdminChatSettings, 0)
	dao.CustomerAdminChatSettings.Ctx(ctx).Where("admin_id in (?)", adminIds).Scan(&settings)

	settingMap := make(map[uint]*entity.CustomerAdminChatSettings)
	for _, s := range settings {
		settingMap[s.AdminId] = s
	}
	for _, message := range res {
		if message.AdminId > 0 && message.Admin != nil {
			setting, exist := settingMap[gconv.Uint(message.AdminId)]
			if exist {
				message.Admin.Setting = setting
			}
		}
	}
	return res
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

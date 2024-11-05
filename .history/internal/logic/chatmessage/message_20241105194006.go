package chatmessage

import (
	"database/sql"
	"gf-chat/internal/consts"
	"gf-chat/internal/dao"
	"gf-chat/internal/model/chat"
	"gf-chat/internal/model/do"
	"gf-chat/internal/model/entity"
	"gf-chat/internal/model/relation"
	"gf-chat/internal/service"
	"time"

	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/gogf/gf/v2/util/grand"
)

func init() {
	service.RegisterChatMessage(&sChatMessage{})
}

type sChatMessage struct {
}

func (s *sChatMessage) GenReqId() string {
	return grand.S(20)
}

func (s *sChatMessage) First(w do.CustomerChatMessages) *entity.CustomerChatMessages {
	message := &entity.CustomerChatMessages{}
	err := dao.CustomerChatMessages.Ctx(gctx.New()).Where(w).Scan(message)
	if err == sql.ErrNoRows {
		return nil
	}
	return message
}

func (s *sChatMessage) SaveRelationOne(msg *relation.CustomerChatMessages) uint {
	result, err := dao.CustomerChatMessages.Ctx(gctx.New()).Save(msg)
	id, err := result.LastInsertId()
	if err != nil {
		return 0
	}
	msg.Id = uint(id)
	return msg.Id
}
func (s *sChatMessage) EntityToRelation(msg *entity.CustomerChatMessages) *relation.CustomerChatMessages {
	return &relation.CustomerChatMessages{
		CustomerChatMessages: *msg,
		Admin:                nil,
		User:                 nil,
	}
}

func (s *sChatMessage) SaveOne(msg *entity.CustomerChatMessages) uint {
	result, err := dao.CustomerChatMessages.Ctx(gctx.New()).Save(msg)
	id, err := result.LastInsertId()
	if err != nil {
		return 0
	}
	msg.Id = uint(id)
	return msg.Id
}

func (s *sChatMessage) ChangeToRead(msgId []uint) (sql.Result, error) {
	ctx := gctx.New()
	query := dao.CustomerChatMessages.Ctx(ctx).Data(map[string]int64{
		"read_at": time.Now().Unix(),
	}).WhereIn("id", msgId).Where("read_at", 0)
	return query.Update()
}

func (s *sChatMessage) GetAdminName(model relation.CustomerChatMessages) string {
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
func (s *sChatMessage) RelationToChat(model relation.CustomerChatMessages) chat.Message {
	username := ""
	if model.User != nil {
		username = model.User.Username
	}
	return chat.Message{
		Id:         model.Id,
		UserId:     model.UserId,
		AdminId:    model.AdminId,
		AdminName:  s.GetAdminName(model),
		Type:       model.Type,
		Content:    model.Content,
		ReceivedAT: model.ReceivedAt,
		Source:     model.Source,
		ReqId:      model.ReqId,
		IsSuccess:  true,
		IsRead:     model.ReadAt != nil,
		Avatar:     s.GetAvatar(model),
		Username:   username,
	}
}
func (s *sChatMessage) GetAvatar(model relation.CustomerChatMessages) string {
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

func (s *sChatMessage) GetModels(lastId int64, w any, size int) []*relation.CustomerChatMessages {
	res := make([]*relation.CustomerChatMessages, 0, 0)
	ctx := gctx.New()
	query := dao.CustomerChatMessages.Ctx(ctx).With(
		relation.CustomerChatMessages{}.Admin,
		// todo
		// 多层关联会用N+1问题
		// relation.CustomerAdmins{}.Setting,
		relation.CustomerChatMessages{}.User,
	).Where(w).OrderDesc("id")
	if size > 0 {
		query = query.Limit(size)
	}
	if lastId > 0 {
		query = query.Where("id < ?", lastId)
	}
	query.Scan(&res)
	adminMaps := make(map[uint]*relation.CustomerAdmins)
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

func (s *sChatMessage) NewNotice(session *entity.CustomerChatSessions, content string) *entity.CustomerChatMessages {
	return &entity.CustomerChatMessages{
		UserId:     session.UserId,
		AdminId:    session.AdminId,
		Type:       consts.MessageTypeNotice,
		Content:    content,
		CustomerId: session.CustomerId,
		SessionId:  session.Id,
		ReqId:      s.GenReqId(),
		Source:     consts.MessageSourceSystem,
	}
}
func (s *sChatMessage) NewOffline(admin *relation.CustomerAdmins) *relation.CustomerChatMessages {
	if admin.Setting != nil && admin.Setting.OfflineContent != "" {
		return &relation.CustomerChatMessages{
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

func (s *sChatMessage) NewWelcome(admin *relation.CustomerAdmins) *relation.CustomerChatMessages {
	if admin.Setting != nil && admin.Setting.WelcomeContent != "" {
		return &relation.CustomerChatMessages{
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

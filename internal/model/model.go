package model

import (
	"gf-chat/api"
	"gf-chat/internal/model/entity"
	"github.com/gogf/gf/v2/frame/g"
)

type CustomerChatAutoRule struct {
	entity.CustomerChatAutoRules
	IsOpen bool     `json:"is_open" orm:"is_open"`
	Scenes []string `json:"scenes" orm:"scenes"`
}

type CustomerChatAutoMessage struct {
	entity.CustomerChatAutoMessages
}

type CustomerAdmin struct {
	g.Meta `orm:"table:customer_admins"`
	entity.CustomerAdmins
	Setting *CustomerAdminChatSetting `orm:"with:admin_id=id"`
}

type CustomerAdminChatSetting struct {
	g.Meta `orm:"table:customer_admin_chat_settings"`
	entity.CustomerAdminChatSettings
	AvatarFile     *CustomerChatFile `orm:"with:id=avatar"`
	BackgroundFile *CustomerChatFile `orm:"with:id=background"`
}

type CustomerChatMessage struct {
	entity.CustomerChatMessages
	Admin *CustomerAdmin `orm:"with:id=admin_id"`
	User  *User          `orm:"with:id=user_id"`
}

type CustomerChatTransfer struct {
	entity.CustomerChatTransfers
	User      *entity.Users          `orm:"with:id=user_id"`
	FormAdmin *entity.CustomerAdmins `orm:"with:id=from_admin_id"`
	ToAdmin   *entity.CustomerAdmins `orm:"with:id=to_admin_id"`
	ToSession *CustomerChatSession   `orm:"with:id=to_session_id"`
}

type CustomerChatSession struct {
	entity.CustomerChatSessions
	User  *User          `orm:"with:id=user_id"`
	Admin *CustomerAdmin `orm:"with:id=admin_id"`
}

type CustomerChatFile struct {
	g.Meta `orm:"table:customer_chat_files"`
	entity.CustomerChatFiles
}

type CustomerChatSetting struct {
	entity.CustomerChatSettings
	Options []api.Option `json:"options"`
}

type User struct {
	g.Meta `orm:"table:users"`
	entity.Users
}

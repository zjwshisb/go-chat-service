package relation

import "gf-chat/internal/model/entity"

type CustomerChatSessions struct {
	entity.CustomerChatSessions
	User  *entity.Users          `orm:"with:id=user_id"`
	Admin *entity.CustomerAdmins `orm:"with:id=admin_id"`
}

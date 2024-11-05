package relation

import (
	"gf-chat/internal/model/entity"
)

type CustomerChatMessages struct {
	entity.CustomerChatMessages
	Admin *CustomerAdmins `orm:"with:id=admin_id"`
	User  *entity.Users   `orm:"with:id=user_id"`
}

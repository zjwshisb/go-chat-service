package relation

import "gf-chat/internal/model/entity"

type CustomerChatTransfer struct {
	entity.CustomerChatTransfers
	User      *entity.Users                `orm:"with:id=user_id"`
	FormAdmin *entity.CustomerAdmins       `orm:"with:id=from_admin_id"`
	ToAdmin   *entity.CustomerAdmins       `orm:"with:id=to_admin_id"`
	ToSession *entity.CustomerChatSessions `orm:"with:id=to_session_id"`
}

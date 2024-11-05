package relation

import (
	"gf-chat/internal/model/entity"
)

type CustomerAdmins struct {
	entity.CustomerAdmins
	Setting *entity.CustomerAdminChatSettings `orm:"with:admin_id=id"`
}

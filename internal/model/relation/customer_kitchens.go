package relation

import "gf-chat/internal/model/entity"

type CustomerKitchen struct {
	entity.CustomerKitchens
	Schools []*entity.Schools `orm:"with:kitchen_id=id"`
}

package relation

import "gf-chat/internal/model/entity"

type Student struct {
	entity.Students
	School entity.Schools `orm:"with:id=school_id"`
}

package user

import (
	"context"
	"database/sql"
	"gf-chat/internal/dao"
	"gf-chat/internal/model/do"
	"gf-chat/internal/model/entity"
	"gf-chat/internal/model/relation"
	"gf-chat/internal/service"

	"github.com/duke-git/lancet/v2/slice"
	"github.com/gogf/gf/v2/os/gctx"
)

func init() {
	service.RegisterUser(&sUser{})
}

type sUser struct {
}

func (s sUser) GetUsers(ctx context.Context, w any) []*entity.Users {
	user := make([]*entity.Users, 0)
	dao.Users.Ctx(ctx).Where(w).
		Scan(&user)
	return user
}

func (s sUser) First(w do.Users) *entity.Users {
	user := &entity.Users{}
	err := dao.Users.Ctx(gctx.New()).Where(w).Scan(user)
	if err == sql.ErrNoRows {
		return nil
	}
	return user
}

func (s sUser) FindByToken(token string) (*entity.Users, *entity.UserApps) {

	return user, userApp
}

func (s sUser) GetRelateStudents(uid int) []relation.Student {
	students := make([]relation.Student, 0)
	dao.Students.Ctx(gctx.New()).
		Where(do.Students{
			UserId:   uid,
			IsCancel: 0,
			IsHide:   0,
		}).WithAll().
		Scan(&students)
	return students
}

func (s sUser) GetStudents(uid int) []entity.Students {
	students := make([]entity.Students, 0)
	dao.Students.Ctx(gctx.New()).
		Where(do.Students{
			UserId:   uid,
			IsCancel: 0,
			IsHide:   0,
		}).
		Scan(&students)
	return students
}

func (s sUser) GetSchools(uid int) []entity.Schools {
	students := s.GetStudents(uid)
	schoolIds := slice.Map(students, func(index int, item entity.Students) int {
		return item.SchoolId
	})
	schools := make([]entity.Schools, 0)
	dao.Schools.Ctx(gctx.New()).WhereIn("id", schoolIds).Scan(&schools)
	return schools
}

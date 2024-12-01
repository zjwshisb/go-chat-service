package cmd

import (
	"context"
	"fmt"
	"gf-chat/internal/model"
	"gf-chat/internal/model/entity"
	"gf-chat/internal/service"
	"golang.org/x/crypto/bcrypt"

	"github.com/gogf/gf/v2/os/gcmd"
)

var Fake = &gcmd.Command{
	Name:        "fake",
	Brief:       "make fake data for testing",
	Description: "make fake data for testing",
	Func: func(ctx context.Context, parser *gcmd.Parser) error {
		for i := range 20 {
			pass, err := bcrypt.GenerateFromPassword([]byte(fmt.Sprintf("admin%d", i)), bcrypt.DefaultCost)
			_, err = service.Admin().Save(ctx, model.CustomerAdmin{
				CustomerAdmins: entity.CustomerAdmins{
					Username:   fmt.Sprintf("admin%d", i),
					Password:   string(pass),
					CustomerId: 1,
				},
			})
			if err != nil {
				panic(err)
			}
			pass, err = bcrypt.GenerateFromPassword([]byte(fmt.Sprintf("user%d", i)), bcrypt.DefaultCost)
			if err != nil {
				panic(err)
			}
			_, err = service.User().Save(ctx, model.User{
				Users: entity.Users{
					Username:   fmt.Sprintf("user%d", i),
					Password:   string(pass),
					CustomerId: 1,
				},
			})
			if err != nil {
				panic(err)
			}
		}
		return nil
	},
}

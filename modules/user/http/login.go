package http

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"ws/models"
	"ws/util"
)

type loginForm struct {
	Username string
	Password string
}

func Login(c *gin.Context) {
	form := &loginForm{}
	err := c.Bind(form)
	if err != nil {
		fmt.Println(err)
	}
	user := &models.User{}
	user.FindByName(form.Username)
	if user.ID !=  0 {
		if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(form.Password)) == nil {
			c.JSON(200, util.RespSuccess(gin.H{
				"token": user.Login(),
			}))
			return
		}
	}
	c.JSON(200, util.RespFail("账号密码错误", 1))
}

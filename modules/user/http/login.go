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
			util.RespSuccess(c, gin.H{
				"token": user.Login(),
			})
			return
		}
	}
	util.RespFail(c, "账号密码错误", 500)
}

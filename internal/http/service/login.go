package service

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"ws/internal/action"
	"ws/internal/models"
	"ws/internal/websocket"
	"ws/util"
)

type loginForm struct {
	Username string `form:"username" binding:"required"`
	Password string `form:"password" binding:"required"`
}

func Login(c *gin.Context) {
	form := &loginForm{}
	err := c.ShouldBind(form)
	if err != nil {
		util.RespValidateFail(c, "表单验证失败")
		return
	}
	user := &models.BackendUser{}
	user.FindByName(form.Username)
	if user.ID !=  0 {
		if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(form.Password)) == nil {
			util.RespSuccess(c, gin.H{
				"token": user.Login(),
			})
			old, exist := websocket.ServiceHub.GetConn(user.ID)
			if exist {
				old.Deliver(action.NewOtherLogin())
			}
			return
		}
	}
	util.RespValidateFail(c, "账号密码错误")
}

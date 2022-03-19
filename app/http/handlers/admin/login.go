package admin

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"strconv"
	"ws/app/http/requests"
	"ws/app/repositories"
	"ws/app/websocket"
)

func Login(c *gin.Context) {
	form := &requests.LoginForm{}
	err := c.ShouldBind(form)
	if err != nil {
		requests.RespValidateFail(c, "表单验证失败")
		return
	}
	user := repositories.AdminRepo.First([]*repositories.Where{
		{
			Filed: "username = ?",
			Value: form.Username,
		},
	}, []string{})
	if user != nil {
		if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(form.Password)) == nil {
			uidStr := strconv.FormatInt(user.GetPrimaryKey(), 10)
			token, _ := requests.CreateToken(uidStr)
			requests.RespSuccess(c, gin.H{
				"token": token,
			})
			websocket.AdminManager.PublishOtherLogin(user)

			return
		}
	}
	requests.RespValidateFail(c, "账号密码错误")
}

package modules

import (
	"ws/modules/service"
	"ws/modules/user"
)

func Setup() {
	user.Setup()
	service.Setup()
}

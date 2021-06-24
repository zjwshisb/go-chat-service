package validators

import (
	"fmt"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"ws/internal/models"
)

func autoMessageTypeValidator(fl validator.FieldLevel) bool {
	if fl.Field().String() == models.TypeText ||
		fl.Field().String() == models.TypeNavigate ||
		fl.Field().String() == models.TypeImage {
		return true
	}
	return false
}

func init() {
	fmt.Println(binding.Validator.Engine())
	v, ok := binding.Validator.Engine().(*validator.Validate)
	if ok {
		_ = v.RegisterValidation("autoMessageType", autoMessageTypeValidator)
	}
}

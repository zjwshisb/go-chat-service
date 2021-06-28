package requests

import (
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"ws/internal/databases"
	"ws/internal/models"
)

func init() {
	v, ok := binding.Validator.Engine().(*validator.Validate)
	if ok {
		_ = v.RegisterValidation("autoMessageType", autoMessageTypeValidator)
		_ = v.RegisterValidation("autoRule", autoRuleValidator)
	}
}

func autoMessageTypeValidator(fl validator.FieldLevel) bool {
	if fl.Field().String() == models.TypeText ||
		fl.Field().String() == models.TypeNavigate{
		if fl.Field().String() == models.TypeNavigate {
			parent := fl.Parent()
			form, _ := parent.Interface().(*AutoMessageForm)
			if form.Url == "" || form.Title == "" {
				return false
			}
		}
		return true
	}
	return false
}
func autoRuleValidator(fl validator.FieldLevel) bool {
	parent := fl.Parent()
	form, _ := parent.Interface().(*AutoRuleForm)
	if form.MatchType != models.MatchTypeAll && form.MatchType != models.MatchTypePart {
		return false
	}
	if form.ReplyType != models.ReplyTypeMessage && form.ReplyType != models.ReplyTypeTransfer {
		return false
	}
	if form.ReplyType == models.ReplyTypeMessage {
		query := databases.Db.Find(&models.AutoMessage{}, form.MessageId)
		if query.RowsAffected == 0 {
			return false
		}
	}
	return true
}


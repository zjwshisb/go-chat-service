package requests

import (
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"ws/app/databases"
	"ws/app/models"
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
		fl.Field().String() == models.TypeNavigate ||
		fl.Field().String() == models.TypeImage {
		if fl.Field().String() == models.TypeNavigate {
			parent := fl.Parent()
			form, _ := parent.Interface().(AutoMessageForm)
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
	form, ok := parent.Interface().(AutoRuleForm)
	if !ok {
		return false
	}
	if len(form.Scenes) > len(models.ScenesOptions) {
		return false
	}
	for _, name := range form.Scenes {
		if name != models.SceneNotAccepted &&
			name != models.SceneAdminOnline &&
			name != models.SceneAdminOffline {
			return false
		}
	}
	if form.MatchType != models.MatchTypeAll && form.MatchType != models.MatchTypePart {
		return false
	}
	if form.ReplyType != models.ReplyTypeMessage && form.ReplyType != models.ReplyTypeTransfer && form.ReplyType !=
		models.ReplyTypeEvent {
		return false
	}
	if form.ReplyType == models.ReplyTypeMessage {
		query := databases.Db.Find(&models.AutoMessage{}, form.MessageId)
		if query.RowsAffected == 0 {
			return false
		}
	}
	if form.ReplyType == models.ReplyTypeEvent{
		if form.Key != models.EventBreak {
			return  false
		}
	}
	return true
}


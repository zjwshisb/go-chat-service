package rule

import (
	"context"
	"gf-chat/internal/service"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/gogf/gf/v2/util/gvalid"
)

func init() {
	gvalid.RegisterRule("auto-message-type", autoMessageTypeRule)
	gvalid.RegisterRule("auto-message-file", autoMessageFile)
}

func autoMessageTypeRule(ctx context.Context, in gvalid.RuleFuncInput) error {
	t := in.Value.String()
	valid := service.ChatMessage().IsTypeValid(t)
	if valid {
		return nil
	}
	return gerror.NewCode(gcode.CodeValidationFailed, "消息类型不正确")
}

func autoMessageFile(ctx context.Context, in gvalid.RuleFuncInput) error {
	form := in.Data.Map()
	types, exist := form["type"]
	if !exist {
		return nil
	}
	if service.ChatMessage().IsFileType(gconv.String(types)) && in.Value.IsEmpty() {
		return gerror.NewCode(gcode.CodeValidationFailed, "请选择文件")
	}
	return nil
}

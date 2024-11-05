package rule

import (
	"context"
	"gf-chat/internal/consts"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/util/gvalid"
)

func init() {
	name := "auto-message-type"
	gvalid.RegisterRule(name, AutoMessageTypeRule)
}

func AutoMessageTypeRule(ctx context.Context, in gvalid.RuleFuncInput) error {
	t := in.Value.String()
	if t == consts.MessageTypeImage || t == consts.MessageTypeText || t == consts.MessageTypeNavigate {
		return nil
	}
	return gerror.New("消息类型不正确")
}

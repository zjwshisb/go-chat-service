package rule

import (
	"context"
	"gf-chat/internal/consts"
	"github.com/duke-git/lancet/v2/slice"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/util/gvalid"
)

func init() {
	gvalid.RegisterRule("auto-rule-match-type", autoRuleMatchTypeRule)
	gvalid.RegisterRule("auto-rule-reply-type", autoRuleReplyTypeRule)
	gvalid.RegisterRule("auto-rule-scene", autoRuleSceneRule)
}

func autoRuleMatchTypeRule(ctx context.Context, in gvalid.RuleFuncInput) error {
	if in.Value.IsEmpty() {
		return nil
	}
	types := []string{
		consts.AutoRuleMatchTypeAll,
		consts.AutoRuleMatchTypePart,
	}
	if slice.Contain(types, in.Value.String()) {
		return nil
	}
	return gerror.NewCode(gcode.CodeValidationFailed, "匹配类型不正确")
}
func autoRuleReplyTypeRule(ctx context.Context, in gvalid.RuleFuncInput) error {
	if in.Value.IsEmpty() {
		return nil
	}
	types := []string{
		consts.AutoRuleReplyTypeMessage,
		consts.AutoRuleReplyTypeTransfer,
	}
	if slice.Contain(types, in.Value.String()) {
		return nil
	}
	return gerror.NewCode(gcode.CodeValidationFailed, "回复类型不正确")
}

func autoRuleSceneRule(ctx context.Context, in gvalid.RuleFuncInput) error {
	if in.Value.IsEmpty() {
		return nil
	}
	types := []string{
		consts.AutoRuleSceneNotAccepted,
		consts.AutoRuleSceneAdminOffline,
		consts.AutoRuleSceneAdminOnline,
	}
	if slice.Contain(types, in.Value.String()) {
		return nil
	}
	return gerror.NewCode(gcode.CodeValidationFailed, "触发场景不正确")
}

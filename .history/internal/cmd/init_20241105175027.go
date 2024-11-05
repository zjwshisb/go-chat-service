package cmd

import (
	"context"
	"database/sql"
	"encoding/json"
	"gf-chat/internal/consts"
	"gf-chat/internal/dao"
	"gf-chat/internal/model/entity"

	_ "github.com/gogf/gf/contrib/drivers/mysql/v2"
	"github.com/gogf/gf/v2/os/gcmd"
)

var (
	options1, _ = json.Marshal([]map[string]string{
		{
			"label": "是",
			"value": "1",
		},
		{
			"label": "否",
			"value": "0",
		},
	})
	options2, _ = json.Marshal([]map[string]string{
		{
			"label": "5分钟",
			"value": "5",
		},
		{
			"label": "10分钟",
			"value": "10",
		},
		{
			"label": "30分钟",
			"value": "30",
		},
		{
			"label": "60分钟",
			"value": "60",
		},
	})
	settings = []entity.CustomerChatSettings{
		{
			Name:        consts.ChatSettingIsAutoTransfer,
			Title:       "自动转接人工客服",
			Value:       "1",
			Options:     string(options1),
			Type:        "select",
			Description: "开启后用户发送任何消息将自动转入到待人工接入列表中，关闭时用户只有发送的消息匹配到转人工的规则才会转入到人工列表中",
		},
		{
			Name:        consts.ChatSettingShowQueue,
			Title:       "显示队列",
			Value:       "0",
			Options:     string(options1),
			Type:        "select",
			Description: "用户等待人工客服接入时，是否显示前面还有多少人在等待",
		},
		{
			Name:        consts.ChatSettingShowRead,
			Title:       "显示已读",
			Value:       "0",
			Options:     string(options1),
			Type:        "select",
			Description: "用户端页面是否显示消息已读/未读",
		},
		{
			Name:        consts.ChatSettingSystemAvatar,
			Title:       "客服系统默认头像",
			Value:       "",
			Options:     "",
			Type:        "image",
			Description: "系统回复消息以及客服没有设置头像时的默认头像",
		},
		{
			Name:        consts.ChatSettingSystemName,
			Title:       "客服系统默认名称",
			Value:       "",
			Options:     "",
			Type:        "text",
			Description: "系统回复消息以及客服没有设置名称时的默认名称",
		},
		{
			Name:        consts.ChatSettingOfflineTmplId,
			Title:       "微信小程序订阅消息模板id",
			Value:       "",
			Options:     "",
			Type:        "select",
			Description: "仅微信小程序，当客服回复用户消息时，如用户不在线，则发送微信小程序订阅消息提醒",
		},
		{
			Name:        consts.ChatSettingOfflineSmsId,
			Title:       "短信提醒模板id",
			Value:       "",
			Options:     "",
			Type:        "select",
			Description: "手动短信提醒的短信模板id",
		},
	}

	rules = []entity.CustomerChatAutoRules{
		{
			Name:      "用户进入客服系统时",
			Match:     consts.AutoRuleMatchEnter,
			MatchType: consts.AutoRuleMatchTypeAll,
			ReplyType: consts.AutoRuleReplyTypeMessage,
			IsSystem:  1,
		},
		{
			Name:      "当转接到人工客服而没有客服在线时(如不设置则继续转接到人工客服)",
			Match:     consts.AutoRuleMatchAdminAllOffLine,
			MatchType: consts.AutoRuleMatchTypeAll,
			ReplyType: consts.AutoRuleReplyTypeMessage,
			IsSystem:  1,
		},
	}

	Init = &gcmd.Command{
		Name:        "init",
		Brief:       "init database",
		Description: "初始化客服数据库",
		Func: func(ctx context.Context, parser *gcmd.Parser) (err error) {
			customers := make([]*entity.Customers, 0, 0)
			dao.Customers.Ctx(ctx).Scan(&customers)
			for _, customer := range customers {
				for _, setting := range settings {
					model := &entity.CustomerChatSettings{}
					err := dao.CustomerChatSettings.Ctx(ctx).
						Where("name", setting.Name).
						Where("customer_id", customer.Id).Scan(model)
					if err == sql.ErrNoRows {
						model = &setting
						model.CustomerId = customer.Id
					} else {
						model.Description = setting.Description
						model.Title = setting.Title
					}
					dao.CustomerChatSettings.Ctx(ctx).Save(model)
				}
				for _, rule := range rules {
					count, _ := dao.CustomerChatAutoRules.Ctx(ctx).
						Where("match", rule.Match).
						Where("is_system", 1).
						Where("customer_id", customer.Id).Count()
					if count == 0 {
						rule.CustomerId = customer.Id
						dao.CustomerChatAutoRules.Ctx(ctx).Save(rule)
					}
				}
			}
			return nil
		},
	}
)

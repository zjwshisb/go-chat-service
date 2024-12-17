package setup

import (
	"gf-chat/api/v1"
	"gf-chat/internal/consts"
	"gf-chat/internal/model"
	"gf-chat/internal/model/do"
	"gf-chat/internal/model/entity"
	"gf-chat/internal/service"

	_ "github.com/gogf/gf/contrib/drivers/mysql/v2"
	"github.com/gogf/gf/v2/os/gctx"
)

var (
	options1 = []v1.Option{
		{
			Label: "是",
			Value: "1",
		},
		{
			Label: "否",
			Value: "0",
		},
	}
	options2 = []v1.Option{
		{
			Label: "5分钟",
			Value: "5",
		},
		{
			Label: "10分钟",
			Value: "10",
		},
		{
			Label: "30分钟",
			Value: "30",
		},
		{
			Label: "60分钟",
			Value: "60",
		},
	}
	settings = []model.CustomerChatSetting{
		{
			CustomerChatSettings: entity.CustomerChatSettings{
				Name:        consts.ChatSettingIsAutoTransfer,
				Title:       "自动转接人工客服",
				Value:       "1",
				Type:        "select",
				Description: "开启后用户发送任何消息将自动转入到待人工接入列表中，关闭时用户只有发送的消息匹配到转人工的规则才会转入到人工列表中",
			},
			Options: options1,
		},
		{
			CustomerChatSettings: entity.CustomerChatSettings{
				Name:        consts.ChatSettingShowQueue,
				Title:       "显示队列",
				Value:       "0",
				Type:        "select",
				Description: "用户等待人工客服接入时，是否显示前面还有多少人在等待",
			},
			Options: options1,
		},
		{
			CustomerChatSettings: entity.CustomerChatSettings{
				Name:        consts.ChatSettingShowRead,
				Title:       "显示已读",
				Value:       "0",
				Type:        "select",
				Description: "用户端页面是否显示消息已读/未读",
			},
			Options: options1,
		},
		{
			CustomerChatSettings: entity.CustomerChatSettings{
				Name:        consts.ChatSettingSystemAvatar,
				Title:       "客服系统默认头像",
				Value:       "",
				Type:        "image",
				Description: "系统回复消息以及客服没有设置头像时的默认头像",
			},
			Options: []v1.Option{},
		},
		{
			CustomerChatSettings: entity.CustomerChatSettings{
				Name:        consts.ChatSettingSystemName,
				Title:       "客服系统默认名称",
				Value:       "",
				Type:        "text",
				Description: "系统回复消息以及客服没有设置名称时的默认名称",
			},
			Options: []v1.Option{},
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
)

func init() {
	service.RegisterSetup(&sSetup{})
}

type sSetup struct {
}

func (s *sSetup) Setup(ctx gctx.Ctx, customerId uint) {
	for _, setting := range settings {
		m, _ := service.ChatSetting().First(ctx, do.CustomerChatSettings{
			Name:       setting.Name,
			CustomerId: setting.CustomerId,
		})
		if m != nil {
			m.CustomerId = customerId
		} else {
			m.Description = setting.Description
			m.Title = setting.Title
		}
		_, err := service.ChatSetting().Save(ctx, m)
		if err != nil {
			panic(err)
		}
	}
	for _, rule := range rules {
		exists, err := service.AutoRule().Exists(ctx, do.CustomerChatAutoRules{
			Match:      rule.Match,
			IsSystem:   1,
			CustomerId: customerId,
		})
		if err != nil {
			panic(err)
		}
		if !exists {
			rule.CustomerId = customerId
			_, err = service.AutoRule().Save(ctx, rules)
			if err != nil {
				panic(err)
			}
		}
	}
}

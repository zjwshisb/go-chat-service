package chatsetting

import (
	"context"
	api "gf-chat/api/v1/backend"
	"gf-chat/internal/consts"
	"gf-chat/internal/dao"
	"gf-chat/internal/model/do"
	"gf-chat/internal/model/entity"
	"gf-chat/internal/service"
	"gf-chat/internal/trait"

	"github.com/gogf/gf/v2/util/gconv"
)

func init() {
	service.RegisterChatSetting(&sChatSetting{
		trait.Curd[entity.CustomerChatSettings]{
			Dao: &dao.CustomerChatSettings,
		},
	})
}

type sChatSetting struct {
	trait.Curd[entity.CustomerChatSettings]
}

func (s *sChatSetting) DefaultAvatarForm(ctx context.Context, customerId uint) (file *api.File, error error) {
	_, err := s.First(ctx, do.CustomerChatSettings{
		CustomerId: customerId,
		Type:       consts.ChatSettingSystemName,
	})
	if err != nil {
		return
	}
	return nil, nil
}

func (s *sChatSetting) GetName(ctx context.Context, customerId uint) (name string, err error) {
	setting, err := s.First(ctx, do.CustomerChatSettings{
		CustomerId: customerId,
		Type:       consts.ChatSettingSystemName,
	})
	if err != nil {
		return
	}
	name = setting.Value
	return
}

func (s *sChatSetting) GetAvatar(ctx context.Context, customerId uint) (name string, err error) {
	setting, err := s.First(ctx, do.CustomerChatSettings{
		CustomerId: customerId,
		Type:       consts.ChatSettingSystemAvatar,
	})
	if err != nil {
		return
	}
	name = setting.Value
	return
}

// GetIsAutoTransferManual 是否自动转接人工客服
func (s *sChatSetting) GetIsAutoTransferManual(ctx context.Context, customerId uint) (b bool, err error) {
	setting, err := s.First(ctx, do.CustomerChatSettings{
		CustomerId: customerId,
		Type:       consts.ChatSettingIsAutoTransfer,
	})
	if err != nil {
		return
	}
	b = gconv.Bool(gconv.Int(setting.Value))
	return
}

package chatsetting

import (
	"database/sql"
	"gf-chat/internal/consts"
	"gf-chat/internal/dao"
	"gf-chat/internal/model"
	"gf-chat/internal/model/entity"
	"gf-chat/internal/service"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/util/gconv"
)

func init() {
	service.RegisterChatSetting(&sChatSetting{})
}

type sChatSetting struct {
}

func (s *sChatSetting) GetSubscribeId(customerId int) string {
	setting := s.First(customerId, consts.ChatSettingOfflineTmplId)
	if setting != nil {
		return setting.Value
	}
	return ""
}

func (s *sChatSetting) First(customerId int, name string) *entity.CustomerChatSettings {
	setting := &entity.CustomerChatSettings{}
	err := dao.CustomerChatSettings.Ctx(gctx.New()).Where("customer_id", customerId).
		Where("name", name).Scan(setting)
	if err == sql.ErrNoRows {
		return nil
	}
	return setting
}

func (s *sChatSetting) DefaultAvatarForm(customerId int) *model.ImageFiled {
	setting := s.First(customerId, consts.ChatSettingSystemAvatar)
	if setting != nil {
		return service.Qiniu().Form(setting.Value)
	}
	return nil
}

func (s *sChatSetting) GetName(customerId int) string {
	setting := s.First(customerId, consts.ChatSettingSystemName)
	return setting.Value
}

func (s *sChatSetting) GetAvatar(customerId int) string {
	setting := s.First(customerId, consts.ChatSettingSystemAvatar)
	if setting != nil {
		return service.Qiniu().Url(setting.Value)
	}
	return ""
}

func (s *sChatSetting) GetSmsCode(customerId int) string {
	setting := s.First(customerId, consts.ChatSettingOfflineSmsId)
	if setting != nil {
		return setting.Value
	}
	return ""
}

// GetIsAutoTransferManual 是否自动转接人工客服
func (s *sChatSetting) GetIsAutoTransferManual(customerId int) bool {
	setting := s.First(customerId, consts.ChatSettingIsAutoTransfer)
	if setting != nil {
		return gconv.Bool(gconv.Int(setting.Value))
	}
	return true
}

package chat

import (
	"strconv"
	"ws/app/databases"
	"ws/app/models"
)



var SettingService = &settingService{
}

type settingService struct {

}



// GetOfflineDuration 客服离线多久自动断开
func (settingService *settingService) GetOfflineDuration(gid int64) int64 {
	setting := &models.ChatSetting{}
	databases.Db.Where("group_id = ?", gid).Where("name = ?", models.MinuteToBreak).First(setting)
	if setting.Id != 0 {
		min, err := strconv.ParseInt(setting.Value, 10, 64)
		if err == nil {
			return  min * 60
		}
	}
	return 5 * 60
}


// GetIsAutoTransferManual 是否自动转接人工客服
func (settingService *settingService) GetIsAutoTransferManual(gid int64) bool {
	setting := &models.ChatSetting{}
	databases.Db.Where("group_id = ?", gid).Where("name = ?", models.IsAutoTransfer).First(setting)
	if setting.Id != 0 {
		min, err := strconv.ParseInt(setting.Value, 10, 64)
		if err == nil {
			return  min > 0
		}
	}
	return true
}


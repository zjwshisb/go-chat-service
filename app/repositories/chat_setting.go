package repositories

import (
	"ws/app/databases"
	"ws/app/models"
)

type ChatSettingRepo struct {
	
}

func (setting *ChatSettingRepo) First(wheres []*Where, orders []string) *models.ChatSetting {
	s := &models.ChatSetting{}
	result := databases.Db.Scopes(AddOrder(orders)).Scopes(AddWhere(wheres)).First(s)
	if result.Error != nil {
		return nil
	}
	return s
}
func (setting *ChatSettingRepo) GetSystemAvatar(groupId int64)  string{
	s := setting.First([]*Where{
		{
			Filed: "type = ?",
			Value: models.SystemAvatar,
		},
		{
			Filed: "group_id = ?",
			Value: groupId,
		},
	}, []string{"id"})
	if s != nil {
		return s.Value
	}
	return ""
}
func (setting *ChatSettingRepo) GetSystemName(groupId int64) string {
	s := setting.First([]*Where{
		{
		Filed: "type = ?",
		Value: models.SystemName,
		},
		{
			Filed: "group_id = ?",
			Value: groupId,
		},
	}, []string{"id"})
	if s != nil {
		return s.Value
	}
	return ""
}



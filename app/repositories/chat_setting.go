package repositories

import (
	"gorm.io/gorm/clause"
	"ws/app/databases"
	"ws/app/models"
)

type chatSettingRepo struct {
	
}
func (setting *chatSettingRepo) Save(s *models.ChatSetting)  {
	databases.Db.Omit(clause.Associations).Save(s)
}
func (setting *chatSettingRepo) Get(wheres []*Where, limit int, load []string, orders []string) []*models.ChatSetting {
	items := make([]*models.ChatSetting, 0 )
	databases.Db.Scopes(AddWhere(wheres)).
		Scopes(AddLoad(load)).
		Scopes(AddOrder(orders)).Limit(limit).Find(&items)
	return items
}

func (setting *chatSettingRepo) First(wheres []*Where, orders []string) *models.ChatSetting {
	s := &models.ChatSetting{}
	result := databases.Db.Scopes(AddOrder(orders)).Scopes(AddWhere(wheres)).First(s)
	if result.Error != nil {
		return nil
	}
	return s
}
func (setting *chatSettingRepo) GetSystemAvatar(groupId int64)  string{
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
func (setting *chatSettingRepo) GetSystemName(groupId int64) string {
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



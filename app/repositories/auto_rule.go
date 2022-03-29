package repositories

import (
	"ws/app/databases"
	"ws/app/models"
)

type autoRuleRepo struct {
	Repository[models.AutoRule]
}

func (repo *autoRuleRepo) Save(model *models.AutoRule) {
	repo.Repository.Save(model)
	for _, scene := range model.Scenes {
		databases.Db.Save(&models.AutoRuleScene{
			Name:   scene.Name,
			RuleId: model.ID,
		})
	}
}

func (repo *autoRuleRepo) DeleteScene(rule *models.AutoRule) int64 {
	result := databases.Db.Where("rule_id = ?", rule.ID).Delete(&models.AutoRuleScene{})
	return result.RowsAffected
}

func (repo *autoRuleRepo) GetWithScenesRuleIds(scene string) []string {
	ids := make([]string, 0)
	databases.Db.Model(&models.AutoRuleScene{}).
		Where("name = ?", scene).
		Pluck("rule_id", &ids)
	return ids
}

// GetAllActiveNormalByGroup 获取所有的启用的普通规则
func (repo *autoRuleRepo) GetAllActiveNormalByGroup(gid int64) []*models.AutoRule {
	return repo.Get([]*Where{
		{
			Filed: "is_system",
			Value: 0,
		},
		{
			Filed: "is_open",
			Value: 1,
		},
		{
			Filed: "group_id",
			Value: gid,
		},
	}, -1, []string{"Message", "Scenes"}, []string{"sort"})
}
func (repo *autoRuleRepo) GetEnterByGroup(gid int64) *models.AutoRule {
	return repo.First([]*Where{
		{
			Filed: "is_system",
			Value: 1,
		},
		{
			Filed: "match",
			Value: models.MatchEnter,
		},
		{
			Filed: "group_id",
			Value: gid,
		},
	}, []string{})
}

// GetAdminAllOffLine 获取转接人工时没有客服在线规则
func (repo *autoRuleRepo) GetAdminAllOffLine(gid int64) *models.AutoRule {
	return repo.Repository.First([]*Where{
		{
			Filed: "is_system = ?",
			Value: 1,
		},
		{
			Filed: "match = ?",
			Value: models.MatchAdminAllOffLine,
		},
		{
			Filed: "gid = ?",
			Value: gid,
		},
	}, []string{})
}

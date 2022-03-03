package repositories

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm/clause"
	"ws/app/databases"
	"ws/app/models"
)

type AutoRuleRepo struct {
}

func (repo *AutoRuleRepo) Update(wheres []*Where, values map[string]interface{}) int64 {
	query := databases.Db.Model(&models.AutoRule{}).Scopes(AddWhere(wheres))
	query.Updates(values)
	return query.RowsAffected
}
func (repo *AutoRuleRepo) Save(model *models.AutoRule)  {
	databases.Db.Omit(clause.Associations).Save(model)
	for _ ,scene := range model.Scenes {
		databases.Db.Save(&models.AutoRuleScene{
			Name:      scene.Name,
			RuleId:    model.ID,
		})
	}
}
func (repo *AutoRuleRepo) Delete(rule *models.AutoRule) int64 {
	result := databases.Db.Delete(rule)
	repo.DeleteScene(rule)
	return result.RowsAffected
}
func (repo *AutoRuleRepo) DeleteScene(rule *models.AutoRule) int64 {
	result := databases.Db.Where("rule_id = ?", rule.ID).Delete(&models.AutoRuleScene{})
	return result.RowsAffected
}

func (repo *AutoRuleRepo) Get(wheres []*Where, limit int, loads []string, orders []string) []*models.AutoRule {
	rules := make([]*models.AutoRule, 0)
	query := databases.Db.
		Limit(limit).
		Scopes(AddWhere(wheres)).Scopes(AddLoad(loads)).Scopes(AddOrder(orders))
	query.Find(&rules)
	return rules
}
func (repo *AutoRuleRepo) GetWithScenesRuleIds(scene string) []string {
	ids := make([]string, 0)
	databases.Db.Model(&models.AutoRuleScene{}).
		Where("name = ?" , scene).
		Pluck("rule_id", &ids)
	return ids
}
func (repo *AutoRuleRepo) Paginate(c *gin.Context, wheres []*Where, load []string, order []string) *Pagination {
	rules := make([]*models.AutoRule, 0)
	databases.Db.
		Scopes(Paginate(c)).
		Scopes(AddWhere(wheres)).
		Scopes(AddLoad(load)).
		Scopes(AddOrder(order)).
		Find(&rules)
	var total int64
	databases.Db.Model(&models.AutoRule{}).
		Scopes(AddWhere(wheres)).
		Count(&total)
	return NewPagination(rules, total)
}

// GetAllActiveNormalByGroup 获取所有的启用的普通规则
func (repo *AutoRuleRepo) GetAllActiveNormalByGroup(gid int64) []*models.AutoRule {
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
	}, -1 , []string{"Message", "Scenes"}, []string{"sort"})
}
func (repo *AutoRuleRepo) GetEnterByGroup(gid int64) *models.AutoRule {
	return repo.First([]*Where{
		{
			Filed: "is_system = ?",
			Value: 1,
		},
		{
			Filed: "match = ?",
			Value: models.MatchEnter,
		},
		{
			Filed: "group_id = ?",
			Value: gid,
		},
	}, []string{})
}

// GetAdminAllOffLine 获取转接人工时没有客服在线规则
func (repo *AutoRuleRepo) GetAdminAllOffLine(gid int64) *models.AutoRule {
	return repo.First([]*Where{
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

func (repo *AutoRuleRepo) First(wheres []*Where, orders []string) *models.AutoRule {
	rule := &models.AutoRule{}
	result := databases.Db.Scopes(AddOrder(orders)).Scopes(AddWhere(wheres)).First(rule)
	if result.Error != nil {
		return nil
	}
	return rule
}


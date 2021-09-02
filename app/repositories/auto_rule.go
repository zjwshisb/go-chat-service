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

func (repo *AutoRuleRepo) Get(wheres []*Where, limit int, loads []string, orders ...string) []*models.AutoRule {
	rules := make([]*models.AutoRule, 0)
	query := databases.Db.
		Limit(limit).
		Scopes(AddWhere(wheres)).Scopes(AddLoad(loads)).Scopes(AddOrder(orders))
	query.Find(&rules)
	return rules
}

func (repo *AutoRuleRepo) Paginate(c *gin.Context, wheres ...*Where) *Pagination {
	rules := make([]*models.AutoRule, 0)
	databases.Db.Order("id desc").
		Scopes(Filter(c, []string{"reply_type"})).
		Scopes(Paginate(c)).
		Scopes(AddWhere(wheres)).
		Preload("Message").
		Find(&rules)
	var total int64
	databases.Db.Model(&models.AutoRule{}).
		Scopes(Filter(c, []string{"reply_type"})).
		Scopes(AddWhere(wheres)).
		Count(&total)
	return NewPagination(rules, total)
}
// 获取所有的启用的普通规则
func (repo *AutoRuleRepo) GetAllActiveNormal() []*models.AutoRule {
	return repo.Get([]*Where{
		{
			Filed: "is_system",
			Value: 0,
		},
		{
			Filed: "is_open",
			Value: 1,
		},
	}, -1 , []string{"Message", "Scenes"}, "sort")
}
func (repo *AutoRuleRepo) GetEnter() *models.AutoRule {
	return repo.First([]*Where{
		{
			Filed: "is_system = ?",
			Value: 1,
		},
		{
			Filed: "match",
			Value: models.MatchEnter,
		},
	})
}
// 获取转接人工时没有客服在线规则
func (repo *AutoRuleRepo) GetAdminAllOffLine() *models.AutoRule {
	return repo.First([]*Where{
		{
			Filed: "is_system",
			Value: 1,
		},
		{
			Filed: "match",
			Value: models.MatchAdminAllOffLine,
		},
	})
}
// get one
func (repo *AutoRuleRepo) First(wheres []*Where, orders ...string) *models.AutoRule {
	rule := &models.AutoRule{}
	result := databases.Db.Scopes(AddOrder(orders)).Scopes(AddWhere(wheres)).First(rule)
	if result.Error != nil {
		return nil
	}
	return rule
}


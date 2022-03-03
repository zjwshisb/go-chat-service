package repositories

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm/clause"
	"time"
	"ws/app/databases"
	"ws/app/models"
)

type ChatSessionRepo struct {
}


func (session *ChatSessionRepo) Create(uid int64, groupId int64, ty int8) *models.ChatSession  {
	s := &models.ChatSession{}
	s.UserId = uid
	s.QueriedAt = time.Now().Unix()
	s.AdminId = 0
	s.Type = ty
	s.GroupId = groupId
	session.Save(s)
	return s
}

func (session *ChatSessionRepo) Save(model *models.ChatSession) int64 {
	result := databases.Db.Omit(clause.Associations).Save(model)
	return result.RowsAffected
}
func (session *ChatSessionRepo) Delete(where []*Where) int64  {
	result := databases.Db.Scopes(AddWhere(where)).Delete(&models.ChatSession{})
	return result.RowsAffected
}
func (session *ChatSessionRepo) First(where []*Where, orders []string) *models.ChatSession  {
	model := &models.ChatSession{}
	databases.Db.Scopes(AddWhere(where)).Scopes(AddOrder(orders)).Find(model)
	if model.Id == 0 {
		return nil
	}
	return model
}

func (session *ChatSessionRepo) GetWaitHandles() []*models.ChatSession  {
	sessions := make([]*models.ChatSession, 0)
	databases.Db.
		Limit(-1).
		Scopes(AddWhere([]*Where{
		{
			Filed: "admin_id = ?",
			Value: 0,
		},
		{
			Filed: "canceled_at = ?",
			Value: 0,
		},
		{
			Filed: "type = ?",
			Value: models.ChatSessionTypeNormal,
		},
	})).
		Preload("Messages", "source = ?" , models.SourceUser).
		Preload("User").
		Find(&sessions)
	return sessions
}
func (session *ChatSessionRepo) Get(wheres []*Where, limit int, load []string, orders []string) []*models.ChatSession {
	sessions := make([]*models.ChatSession, 0)
	databases.Db.
		Limit(limit).
		Scopes(AddOrder(orders)).
		Scopes(AddWhere(wheres)).
		Scopes(AddLoad(load)).
		Find(&sessions)
	return sessions
}

// FirstActiveByUser 获取有效会话
func (session *ChatSessionRepo) FirstActiveByUser(uid int64, adminId int64) *models.ChatSession {
	s := session.First([]*Where{
		{
			Filed: "user_id = ?",
			Value: uid,
		},
		{
			Filed: "admin_id = ?",
			Value: adminId,
		},
		{
			Filed: "broke_at = ? ",
			Value: 0,
		},
		{
			Filed: "canceled_at = ?",
			Value: 0,
		},
	}, []string{"id desc"})
	return s
}

func (session *ChatSessionRepo) Paginate(c *gin.Context, wheres []*Where, load []string, orders []string) *Pagination {
	sessions := make([]*models.ChatSession, 0)
	databases.Db.Order("id desc").
		Scopes(Paginate(c)).
		Scopes(AddLoad(load)).
		Scopes(AddOrder(orders)).
		Scopes(AddWhere(wheres)).
		Find(&sessions)
	var total int64
	databases.Db.Model(&models.ChatSession{}).
		Scopes(AddWhere(wheres)).
		Count(&total)
	return NewPagination(sessions, total)
}


package repositories

import (
	"ws/app/databases"
	"ws/app/models"
)
type TransferRepo struct {
}

func (repo *TransferRepo) Save(transfer *models.ChatTransfer) int64 {
	result := databases.Db.Save(transfer)
	return result.RowsAffected
}

func (repo *TransferRepo) First(where []*Where, orders ...string) *models.ChatTransfer  {
	model := &models.ChatTransfer{}
	result := databases.Db.Scopes(AddWhere(where)).Scopes(AddOrder(orders)).First(model)
	if result.Error != nil {
		return nil
	}
	return model
}

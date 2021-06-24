package backend

import (
	"github.com/gin-gonic/gin"
	"ws/internal/databases"
	"ws/internal/models"
	"ws/internal/util"
)

func UpdateSetting(c *gin.Context) {

}

func GetSettings(c *gin.Context) {
	var settings []models.Setting
	databases.Db.Find(&settings)
	util.RespSuccess(c, settings)
}
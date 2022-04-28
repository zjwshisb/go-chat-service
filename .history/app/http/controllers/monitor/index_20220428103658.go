package monitor

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Monitor(c *gin.Context) {
	c.HTML(http.StatusOK, "monitor.tmpl", gin.H{})
}

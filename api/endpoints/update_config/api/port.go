package api

import (
	"eigen_db/cfg"
	"net/http"

	"github.com/gin-gonic/gin"
)

type updatePortBody struct {
	UpdatedPort int `json:"updatedPort" binding:"required"`
}

func UpdatePort(config cfg.IConfig) func(*gin.Context) {
	return func(c *gin.Context) {
		var body updatePortBody
		if err := c.ShouldBindJSON(&body); err != nil {
			c.Status(http.StatusBadRequest)
			return
		}

		_ = config.SetAPIPort(body.UpdatedPort) // handle error
		c.String(http.StatusOK, "API port updated. Please restart the database for it to take effect.")
	}
}

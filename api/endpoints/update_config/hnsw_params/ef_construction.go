package hnsw_params

import (
	"eigen_db/cfg"
	"net/http"

	"github.com/gin-gonic/gin"
)

type updateEfConstBody struct {
	UpdatedEfConst int `json:"updatedEfConst" binding:"required"`
}

func UpdateEfConstruction(config cfg.IConfig) func(*gin.Context) {
	return func(c *gin.Context) {
		var body updateEfConstBody
		if err := c.ShouldBindJSON(&body); err != nil {
			c.Status(http.StatusBadRequest)
			return
		}

		config.SetHNSWParamsEfConstruction(body.UpdatedEfConst)
		c.String(http.StatusOK, "EF Construction paramater updated. Please restart the database for it to take effect.")
	}
}

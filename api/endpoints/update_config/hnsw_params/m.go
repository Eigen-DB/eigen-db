package hnsw_params

import (
	"eigen_db/cfg"
	"net/http"

	"github.com/gin-gonic/gin"
)

type updateMBody struct {
	UpdatedM int `json:"updatedM" binding:"required"`
}

func UpdateM(config cfg.IConfig) func(*gin.Context) {
	return func(c *gin.Context) {
		var body updateMBody
		if err := c.ShouldBindJSON(&body); err != nil {
			c.Status(http.StatusBadRequest)
			return
		}

		_ = config.SetHNSWParamsM(body.UpdatedM) // handle error
		c.String(http.StatusOK, "M paramater updated. Please restart the database for it to take effect.")
	}
}

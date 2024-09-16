package hnsw_params

import (
	"eigen_db/cfg"
	"net/http"

	"github.com/gin-gonic/gin"
)

type updateSpaceSizeBody struct {
	UpdatedSize uint32 `json:"updatedSize" binding:"required"`
}

func UpdateSpaceSize(config cfg.IConfig) func(*gin.Context) {
	return func(c *gin.Context) {
		var body updateSpaceSizeBody
		if err := c.ShouldBindJSON(&body); err != nil {
			c.Status(http.StatusBadRequest)
			return
		}

		_ = config.SetHNSWParamsSpaceSize(body.UpdatedSize) // handle error
		c.String(http.StatusOK, "Vector space size updated. Please restart the database for it to take effect.")
	}
}

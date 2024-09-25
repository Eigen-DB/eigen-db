package hnsw_params

import (
	"eigen_db/api/utils"
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
		if err := utils.ValidateBody(c, &body); err != nil {
			return
		}

		config.SetHNSWParamsSpaceSize(body.UpdatedSize)
		utils.SendResponse(
			c,
			http.StatusOK,
			"Vector space size updated. Please restart the database for it to take effect.",
			nil,
			nil,
		)
	}
}

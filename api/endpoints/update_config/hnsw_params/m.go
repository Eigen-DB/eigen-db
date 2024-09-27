package hnsw_params

import (
	"eigen_db/api/utils"
	"eigen_db/cfg"
	"net/http"

	"github.com/gin-gonic/gin"
)

type updateMBody struct {
	UpdatedM int `json:"updatedM" binding:"required,gt=0"`
}

func UpdateM(config *cfg.Config) func(*gin.Context) {
	return func(c *gin.Context) {
		var body updateMBody
		if err := utils.ValidateBody(c, &body); err != nil {
			return
		}

		config.SetHNSWParamsM(body.UpdatedM)
		utils.SendResponse(
			c,
			http.StatusOK,
			"M paramater updated. Please restart the database for it to take effect.",
			nil,
			nil,
		)
	}
}

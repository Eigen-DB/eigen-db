package api

import (
	"eigen_db/api/utils"
	"eigen_db/cfg"
	"net/http"

	"github.com/gin-gonic/gin"
)

type updatePortBody struct {
	UpdatedPort int `json:"updatedPort" binding:"required"`
}

func UpdatePort(config *cfg.Config) func(*gin.Context) {
	return func(c *gin.Context) {
		var body updatePortBody
		if err := utils.ValidateBody(c, &body); err != nil {
			return
		}

		config.SetAPIPort(body.UpdatedPort)
		utils.SendResponse(
			c,
			http.StatusOK,
			"API port updated. Please restart the database for it to take effect.",
			nil,
			nil,
		)
	}
}

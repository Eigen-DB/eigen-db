package persistence

import (
	"eigen_db/api/utils"
	"eigen_db/cfg"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type updateTimeIntervalBody struct {
	UpdatedValueSecs float32 `json:"updatedValueSecs" binding:"required,gt=0"`
}

func UpdateTimeInterval(config cfg.IConfig) func(*gin.Context) {
	return func(c *gin.Context) {
		var body updateTimeIntervalBody
		if err := utils.ValidateBody(c, &body); err != nil {
			return
		}

		config.SetPersistenceTimeInterval(time.Duration(body.UpdatedValueSecs * 1.0e+9))
		utils.SendResponse(
			c,
			http.StatusOK,
			"Time interval updated.",
			nil,
			nil,
		)
	}
}

package persistence

import (
	"eigen_db/cfg"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type updateTimeIntervalBody struct {
	UpdatedValueSecs float32 `json:"updatedValueSecs" binding:"required"`
}

func UpdateTimeInterval(config cfg.IConfig) func(*gin.Context) {
	return func(c *gin.Context) {
		var body updateTimeIntervalBody
		if err := c.ShouldBindJSON(&body); err != nil {
			c.Status(http.StatusBadRequest)
			return
		}

		_ = config.SetPersistenceTimeInterval(time.Duration(body.UpdatedValueSecs * 1.0e+9)) // handle error
		c.String(http.StatusOK, "Time interval updated.")
	}
}

package persistence

import (
	"eigen_db/api/utils"
	"eigen_db/cfg"
	"eigen_db/constants"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type updateTimeIntervalBody struct {
	UpdatedValueSecs float32 `json:"updatedValueSecs" binding:"required,gt=0"`
}

func UpdateTimeInterval(c *gin.Context) {
	var body updateTimeIntervalBody
	if err := utils.ValidateBody(c, &body); err != nil {
		return
	}

	config := cfg.GetConfig()
	config.SetPersistenceTimeInterval(time.Duration(body.UpdatedValueSecs * 1.0e+9))
	err := config.WriteToDisk(constants.CONFIG_PATH)
	if err != nil {
		utils.SendResponse(
			c,
			http.StatusInternalServerError,
			"An error occured.",
			nil,
			utils.CreateError("ERROR_UPDATING_PERSISTENCE_TIME_INTERVAL", fmt.Sprintf("Error: %s", err.Error())),
		)
		return
	}

	utils.SendResponse(
		c,
		http.StatusOK,
		"Time interval updated.",
		nil,
		nil,
	)
}

package api

import (
	"eigen_db/api/utils"
	"eigen_db/cfg"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type updatePortBody struct {
	UpdatedPort int `json:"updatedPort" binding:"required"`
}

func UpdatePort(c *gin.Context) {
	var body updatePortBody
	if err := utils.ValidateBody(c, &body); err != nil {
		return
	}

	config := cfg.GetConfig()
	err := config.SetAPIPort(body.UpdatedPort)
	if err != nil {
		utils.SendResponse(
			c,
			http.StatusInternalServerError,
			"An error occured.",
			nil,
			utils.CreateError("ERROR_UPDATING_API_PORT", fmt.Sprintf("Error: %s", err.Error())),
		)
		return
	}

	utils.SendResponse(
		c,
		http.StatusOK,
		"API port updated. Please restart the database for it to take effect.",
		nil,
		nil,
	)
}

package api

import (
	"eigen_db/api/utils"
	"eigen_db/cfg"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type updateAddressBody struct {
	UpdatedAddress string `json:"updatedAddress" binding:"required"`
}

func UpdateAddress(c *gin.Context) {
	var body updateAddressBody
	if err := utils.ValidateBody(c, &body); err != nil {
		return
	}

	config := cfg.GetConfig()
	err := config.SetAPIAddress(body.UpdatedAddress)
	if err != nil {
		utils.SendResponse(
			c,
			http.StatusInternalServerError,
			"An error occured.",
			nil,
			utils.CreateError("ERROR_UPDATING_API_ADDRESS", fmt.Sprintf("Error: %s", err.Error())),
		)
		return
	}

	utils.SendResponse(
		c,
		http.StatusOK,
		"API address updated. Please restart the database for it to take effect.",
		nil,
		nil,
	)
}

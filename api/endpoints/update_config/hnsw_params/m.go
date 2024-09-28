package hnsw_params

import (
	"eigen_db/api/utils"
	"eigen_db/cfg"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type updateMBody struct {
	UpdatedM int `json:"updatedM" binding:"required,gt=0"`
}

func UpdateM(c *gin.Context) {
	var body updateMBody
	if err := utils.ValidateBody(c, &body); err != nil {
		return
	}

	config := cfg.GetConfig()
	err := config.SetHNSWParamsM(body.UpdatedM)
	if err != nil {
		utils.SendResponse(
			c,
			http.StatusInternalServerError,
			"An error occured.",
			nil,
			utils.CreateError("ERROR_UPDATING_M", fmt.Sprintf("Error: %s", err.Error())),
		)
		return
	}

	utils.SendResponse(
		c,
		http.StatusOK,
		"M paramater updated. Please restart the database for it to take effect.",
		nil,
		nil,
	)
}

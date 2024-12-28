package hnsw_params

import (
	"eigen_db/api/utils"
	"eigen_db/cfg"
	"eigen_db/constants"
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
	if err := config.SetM(body.UpdatedM); err != nil {
		utils.SendResponse(
			c,
			http.StatusBadRequest,
			"Invalid M parameter.",
			nil,
			utils.CreateError("INVALID_M", err.Error()),
		)
		return
	}
	if err := config.WriteToDisk(constants.CONFIG_PATH); err != nil {
		utils.SendResponse(
			c,
			http.StatusInternalServerError,
			"An error occured.",
			nil,
			utils.CreateError("ERROR_UPDATING_M", err.Error()),
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

package hnsw_params

import (
	"eigen_db/api/utils"
	"eigen_db/cfg"
	"eigen_db/constants"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type updateSpaceSizeBody struct {
	UpdatedSize uint32 `json:"updatedSize" binding:"required"`
}

func UpdateSpaceSize(c *gin.Context) {
	var body updateSpaceSizeBody
	if err := utils.ValidateBody(c, &body); err != nil {
		return
	}

	config := cfg.GetConfig()
	config.SetSpaceSize(body.UpdatedSize)
	err := config.WriteToDisk(constants.CONFIG_PATH)
	if err != nil {
		utils.SendResponse(
			c,
			http.StatusInternalServerError,
			"An error occured.",
			nil,
			utils.CreateError("ERROR_UPDATING_SPACE_SIZE", fmt.Sprintf("Error: %s", err.Error())),
		)
		return
	}

	utils.SendResponse(
		c,
		http.StatusOK,
		"Vector space size updated. Please restart the database for it to take effect.",
		nil,
		nil,
	)
}

package hnsw_params

import (
	"eigen_db/api/utils"
	"eigen_db/cfg"
	"net/http"

	"github.com/gin-gonic/gin"
)

type updateEfConstBody struct {
	UpdatedEfConst int `json:"updatedEfConst" binding:"required,gt=0"`
}

func UpdateEfConstruction(c *gin.Context) {
	var body updateEfConstBody
	if err := utils.ValidateBody(c, &body); err != nil {
		return
	}

	config := cfg.GetConfig()
	config.SetHNSWParamsEfConstruction(body.UpdatedEfConst)
	utils.SendResponse(
		c,
		http.StatusOK,
		"EF Construction paramater updated. Please restart the database for it to take effect.",
		nil,
		nil,
	)
}

package index_config

import (
	"eigen_db/api/utils"
	"eigen_db/types"
	"net/http"

	"github.com/gin-gonic/gin"
)

type updateSimMetricBody struct {
	UpdatedMetric types.SimMetric `json:"updatedMetric" binding:"required"`
}

func UpdateSimilarityMetric(c *gin.Context) {
	utils.SendResponse(
		c,
		http.StatusInternalServerError,
		"Not implemented yet.",
		nil,
		utils.CreateError("NOT_IMPLEMENTED", "This endpoint is not implemented yet."),
	)

	// ISSUE: since PQ computes a distance table, updating the similarity metric and restarting the DB *might* not update the distance table with the new metric.
	// look into this in the Faiss source code...

	/*
		var body updateSimMetricBody
		if err := utils.ValidateBody(c, &body); err != nil {
			return
		}

		config := cfg.GetConfig()
		if err := config.SetSimilarityMetric(body.UpdatedMetric); err != nil {
			utils.SendResponse(
				c,
				http.StatusBadRequest,
				"Something went wrong when trying to update the similarity metric.",
				nil,
				utils.CreateError("INVALID_SIMILARITY_METRIC", err.Error()),
			)
			return
		}
		if err := config.WriteToDisk(constants.CONFIG_PATH); err != nil {
			utils.SendResponse(
				c,
				http.StatusInternalServerError,
				"An error occured.",
				nil,
				utils.CreateError("ERROR_UPDATING_SIMILARITY_METRIC", err.Error()),
			)
			return
		}

		utils.SendResponse(
			c,
			http.StatusOK,
			"Vector similarity metric updated. Please restart the database for it to take effect.",
			nil,
			nil,
		)*/
}

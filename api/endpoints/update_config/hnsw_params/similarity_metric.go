package hnsw_params

import (
	"eigen_db/api/utils"
	"eigen_db/cfg"
	"eigen_db/types"
	"net/http"

	"github.com/gin-gonic/gin"
)

type updateSimMetricBody struct {
	UpdatedMetric types.SimilarityMetric `json:"updatedMetric" binding:"required"`
}

func UpdateSimilarityMetric(config cfg.IConfig) func(*gin.Context) {
	return func(c *gin.Context) {
		var body updateSimMetricBody
		if err := utils.ValidateBody(c, &body); err != nil {
			return
		}

		metric, err := types.ParseSimilarityMetric(body.UpdatedMetric)
		if err != nil {
			utils.SendResponse(
				c,
				http.StatusBadRequest,
				"Something went wrong when trying to update the similarity metric.",
				nil,
				utils.CreateError("INVALID_SIMILARITY_METRIC", err.Error()),
			)
			return
		}

		config.SetHNSWParamsSimilarityMetric(metric)
		utils.SendResponse(
			c,
			http.StatusOK,
			"Vector similarity metric updated. Please restart the database for it to take effect.",
			nil,
			nil,
		)
	}
}

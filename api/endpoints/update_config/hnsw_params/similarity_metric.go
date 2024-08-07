package hnsw_params

import (
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
		if err := c.ShouldBindJSON(&body); err != nil {
			c.Status(http.StatusBadRequest)
			return
		}

		metric, err := types.ParseSimilarityMetric(body.UpdatedMetric)
		if err != nil {
			c.String(http.StatusBadRequest, err.Error())
			return
		}

		config.SetHNSWParamsSimilarityMetric(metric)
		c.String(http.StatusOK, "Vector similarity metric updated. Please restart the database for it to take effect.")
	}
}

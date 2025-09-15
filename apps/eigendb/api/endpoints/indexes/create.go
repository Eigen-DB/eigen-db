package indexes

import (
	"eigen_db/api/utils"
	"eigen_db/index_mgr"
	"eigen_db/types"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type createRequestBody struct {
	Dimensions int             `json:"dimensions" binding:"required,gt=1"`
	Metric     types.SimMetric `json:"metric" binding:"required"`
}

func Create(c *gin.Context) {
	var body createRequestBody
	if err := utils.ValidateBody(c, &body); err != nil {
		return
	}

	if err := body.Metric.Validate(); err != nil {
		utils.SendResponse(
			c,
			http.StatusBadRequest,
			"Invalid similarity metric provided.",
			nil,
			utils.CreateError("INVALID_SIMILARITY_METRIC", err.Error()),
		)
		return
	}

	indexName := c.Param("index")
	if err := index_mgr.GetIndexMgr().CreateIndex(
		indexName,
		body.Dimensions,
		body.Metric,
	); err != nil {
		utils.SendResponse(
			c,
			http.StatusInternalServerError,
			"An error occured while creating the index.",
			nil,
			utils.CreateError("INDEX_NOT_CREATED", err.Error()),
		)
		return
	}

	utils.SendResponse(
		c,
		http.StatusOK,
		fmt.Sprintf("Index '%s' created successfully.", indexName),
		nil,
		nil,
	)
}

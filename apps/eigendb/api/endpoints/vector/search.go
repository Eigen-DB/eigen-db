package vector

import (
	"eigen_db/api/utils"
	"eigen_db/vector_io"
	"net/http"

	"github.com/gin-gonic/gin"
)

type searchRequestBody struct {
	QueryVector vector_io.Embedding `json:"queryVector" binding:"required"`
	K           int64               `json:"k" binding:"required,gt=0"`
}

func Search(c *gin.Context) {
	var body searchRequestBody
	if err := utils.ValidateBody(c, &body); err != nil {
		return
	}

	nnIds, err := vector_io.GetMemoryIndex().Search(&body.QueryVector, body.K)
	if err != nil {
		utils.SendResponse(
			c,
			http.StatusBadRequest,
			"An error occured during the similarity search.",
			nil,
			utils.CreateError("SIMILARITY_SEARCH_ERROR", err.Error()),
		)
		return
	}

	utils.SendResponse(
		c,
		http.StatusOK,
		"Similarity search successfully performed.",
		map[string]any{
			"nearest_neighbor_ids": nnIds,
		},
		nil,
	)
}

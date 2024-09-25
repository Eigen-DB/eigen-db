package vector

import (
	"eigen_db/api/utils"
	"eigen_db/vector_io"
	"net/http"

	t "eigen_db/types"

	"github.com/gin-gonic/gin"
)

type searchRequestBody struct {
	QueryVectorId t.VectorId `json:"queryVectorId" binding:"required"`
	K             int        `json:"k" binding:"required,gt=0"`
}

func Search(searcher vector_io.IVectorSearcher) func(*gin.Context) {
	return func(c *gin.Context) {
		var body searchRequestBody
		if err := utils.ValidateBody(c, &body); err != nil {
			return
		}

		nnIds, err := searcher.SimilaritySearch(body.QueryVectorId, body.K)
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
}

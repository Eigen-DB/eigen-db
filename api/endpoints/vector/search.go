package vector

import (
	"eigen_db/vector_io"
	"encoding/json"
	"net/http"

	t "eigen_db/types"

	"github.com/gin-gonic/gin"
)

type searchRequestBody struct {
	QueryVectorId t.VectorId `json:"queryVectorId" binding:"required"`
	K             uint32     `json:"k" binding:"required"`
}

func Search(searcher vector_io.IVectorSearcher) func(*gin.Context) {
	return func(c *gin.Context) {
		var body searchRequestBody
		if err := c.ShouldBindJSON(&body); err != nil {
			c.Status(http.StatusBadRequest)
			return
		}

		nnIds, err := searcher.SimilaritySearch(body.QueryVectorId, body.K)
		if err != nil {
			c.String(http.StatusBadRequest, err.Error())
		}

		jsonResponse, err := json.Marshal(nnIds)
		if err != nil {
			c.String(http.StatusBadRequest, err.Error())
		}

		c.String(http.StatusOK, string(jsonResponse))
	}
}

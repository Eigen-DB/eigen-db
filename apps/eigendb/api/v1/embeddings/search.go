package embeddings

import (
	"eigen_db/api/utils"
	"eigen_db/index_mgr"
	"eigen_db/types"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type searchRequestBody struct {
	QueryVector types.EmbeddingData `json:"queryVector" binding:"required"`
	K           int64               `json:"k" binding:"required,gt=0"`
}

func Search(c *gin.Context) {
	var body searchRequestBody
	if err := utils.ValidateBody(c, &body); err != nil {
		return
	}

	indexName := c.Param("index")
	idx, err := index_mgr.GetIndexMgr().GetIndex(indexName)
	if err != nil {
		utils.SendResponse(
			c,
			http.StatusInternalServerError,
			"An error occured while fetching the index.",
			nil,
			utils.CreateError("INDEX_NOT_FETCHED", err.Error()),
		)
		return
	}

	nn, err := idx.Search(body.QueryVector, body.K)
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

	// converting EmbId to string in nn map for coventional JSON
	nnFormatted := make(map[string]map[string]any, len(nn))
	for k, v := range nn {
		nnFormatted[fmt.Sprint(k)] = v
	}

	utils.SendResponse(
		c,
		http.StatusOK,
		"Similarity search successfully performed.",
		map[string]any{
			"nearest_neighbors": nnFormatted,
		},
		nil,
	)
}

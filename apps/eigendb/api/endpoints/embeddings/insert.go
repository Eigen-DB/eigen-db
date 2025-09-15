package embeddings

import (
	"eigen_db/api/utils"
	"eigen_db/index"
	"eigen_db/index_mgr"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type insertRequestBody struct {
	Embeddings []index.Embedding `json:"embeddings" binding:"required"`
}

func Insert(c *gin.Context) {
	var body insertRequestBody
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

	embeddingsInserted := 0
	errors := make([]string, 0)
	for _, embedding := range body.Embeddings {
		v := index.EmbeddingFactory(embedding.Data, embedding.Metadata, embedding.Id)
		if err := idx.Insert(v); err != nil {
			errors = append(errors, fmt.Sprintf("embedding with ID %d was not inserted - %s", embedding.Id, err.Error()))
		} else {
			embeddingsInserted++
		}
	}

	if len(errors) != 0 {
		utils.SendResponse(
			c,
			http.StatusInternalServerError,
			fmt.Sprintf("%d/%d embeddings successfully inserted.", embeddingsInserted, len(body.Embeddings)),
			nil,
			utils.CreateError("EMBEDDINGS_SKIPPED", errors),
		)
	} else {
		utils.SendResponse(
			c,
			http.StatusOK,
			fmt.Sprintf("%d/%d embeddings successfully inserted.", embeddingsInserted, len(body.Embeddings)),
			nil,
			nil,
		)
	}
}

package embeddings

import (
	"eigen_db/api/utils"
	"eigen_db/vector_io"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type insertRequestBody struct {
	Embeddings []vector_io.Embedding `json:"embeddings" binding:"required"`
}

func Insert(c *gin.Context) {
	var body insertRequestBody
	if err := utils.ValidateBody(c, &body); err != nil {
		return
	}

	embeddingsInserted := 0
	errors := make([]string, 0)
	for _, embedding := range body.Embeddings {
		v, err := vector_io.EmbeddingFactory(embedding.Data, embedding.Metadata, embedding.Id)
		if err != nil {
			errors = append(errors, fmt.Sprintf("embedding with ID %d was not inserted - %s", embedding.Id, err.Error()))
			continue
		}

		if err := vector_io.GetMemoryIndex().Insert(v); err != nil {
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

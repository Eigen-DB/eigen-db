package vector

import (
	"eigen_db/api/utils"
	"eigen_db/vector_io"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type upsertRequestBody struct {
	Embeddings []vector_io.Embedding `json:"embeddings" binding:"required"`
}

func Upsert(c *gin.Context) {
	var body upsertRequestBody
	if err := utils.ValidateBody(c, &body); err != nil {
		return
	}

	embeddingsUpserted := 0
	errors := make([]string, 0)
	for _, embedding := range body.Embeddings {
		v, err := vector_io.EmbeddingFactory(embedding.Data, embedding.Metadata, embedding.Id)
		if err != nil {
			errors = append(errors, fmt.Sprintf("embedding with ID %d was not upserted - %s", embedding.Id, err.Error()))
			continue
		}

		if err := vector_io.GetMemoryIndex().Upsert(v); err != nil {
			errors = append(errors, fmt.Sprintf("embedding with ID %d was not upserted - %s", embedding.Id, err.Error()))
		} else {
			embeddingsUpserted++
		}
	}

	if len(errors) != 0 {
		utils.SendResponse(
			c,
			http.StatusInternalServerError,
			fmt.Sprintf("%d/%d embeddings successfully upserted.", embeddingsUpserted, len(body.Embeddings)),
			nil,
			utils.CreateError("EMBEDDINGS_SKIPPED", errors),
		)
	} else {
		utils.SendResponse(
			c,
			http.StatusOK,
			fmt.Sprintf("%d/%d embeddings successfully upserted.", embeddingsUpserted, len(body.Embeddings)),
			nil,
			nil,
		)
	}
}

package vector

import (
	"eigen_db/api/utils"
	"eigen_db/vector_io"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type bulkInsertRequestBody struct {
	Vectors []vector_io.Embedding `json:"embeddings" binding:"required"`
}

func BulkInsert(c *gin.Context) {
	var body bulkInsertRequestBody
	if err := utils.ValidateBody(c, &body); err != nil {
		return
	}

	vectorsInserted := 0
	errors := make([]string, 0)
	for _, vector := range body.Vectors {
		v, err := vector_io.EmbeddingFactory(vector.Data, vector.Metadata, vector.Id)
		if err != nil {
			errors = append(errors, fmt.Sprintf("vector with ID %d was skipped - %s", vector.Id, err.Error()))
			continue
		}

		if err := vector_io.GetMemoryIndex().Insert(v); err != nil {
			errors = append(errors, fmt.Sprintf("vector with ID %d was skipped - %s", vector.Id, err.Error()))
		} else {
			vectorsInserted++
		}
	}

	if len(errors) != 0 {
		utils.SendResponse(
			c,
			http.StatusInternalServerError,
			fmt.Sprintf("%d/%d vectors successfully inserted.", vectorsInserted, len(body.Vectors)),
			nil,
			utils.CreateError("VECTORS_SKIPPED", errors),
		)
	} else {
		utils.SendResponse(
			c,
			http.StatusOK,
			fmt.Sprintf("%d/%d vectors successfully inserted.", vectorsInserted, len(body.Vectors)),
			nil,
			nil,
		)
	}
}

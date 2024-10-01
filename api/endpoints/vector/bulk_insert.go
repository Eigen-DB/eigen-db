package vector

import (
	"eigen_db/api/utils"
	t "eigen_db/types"
	"eigen_db/vector_io"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type bulkInsertRequestBody struct {
	Embeddings [][]t.VectorComponent `json:"embeddings" binding:"required"`
}

func BulkInsert(c *gin.Context) {
	var body bulkInsertRequestBody
	if err := utils.ValidateBody(c, &body); err != nil {
		return
	}

	vectorsInserted := 0
	errors := make([]string, 0)
	for i, components := range body.Embeddings {
		v, err := vector_io.NewVector(components)
		if err != nil {
			errors = append(errors, fmt.Sprintf("vector %d was skipped - %s", i+1, err.Error()))
			continue
		}

		if err := vector_io.InsertVector(v); err != nil {
			errors = append(errors, fmt.Sprintf("vector %d was skipped - %s", i+1, err.Error()))
		} else {
			vectorsInserted++
		}
	}

	if len(errors) != 0 {
		utils.SendResponse(
			c,
			http.StatusInternalServerError,
			fmt.Sprintf("%d/%d vectors successfully inserted.", vectorsInserted, len(body.Embeddings)),
			nil,
			utils.CreateError("VECTORS_SKIPPED", errors),
		)
	} else {
		utils.SendResponse(
			c,
			http.StatusOK,
			fmt.Sprintf("%d/%d vectors successfully inserted.", vectorsInserted, len(body.Embeddings)),
			nil,
			nil,
		)
	}
}

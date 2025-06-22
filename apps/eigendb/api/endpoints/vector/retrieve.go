package vector

import (
	"eigen_db/api/utils"
	"eigen_db/types"
	"eigen_db/vector_io"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type retrieveRequestBody struct {
	Ids []types.EmbId `json:"ids" binding:"required"`
}

func Retrieve(c *gin.Context) {
	var body retrieveRequestBody
	if err := utils.ValidateBody(c, &body); err != nil {
		return
	}

	embeddings := make([]map[string]any, 0, len(body.Ids))
	errors := make([]string, 0)
	for _, id := range body.Ids {
		embedding, err := vector_io.GetMemoryIndex().Get(id)
		if err != nil {
			errors = append(errors, fmt.Sprintf("embedding with ID %d was not retrieved - %s", id, err.Error()))
			continue
		} else {
			embeddings = append(embeddings, map[string]any{
				"id":       embedding.Id,
				"data":     embedding.Data,
				"metadata": embedding.Metadata,
			})
		}
	}

	if len(errors) != 0 {
		utils.SendResponse(
			c,
			http.StatusInternalServerError,
			fmt.Sprintf("%d/%d embeddings successfully retrieved.", len(embeddings), len(body.Ids)),
			map[string]any{"embeddings": embeddings},
			utils.CreateError("EMBEDDINGS_SKIPPED", errors),
		)
	} else {
		utils.SendResponse(
			c,
			http.StatusOK,
			fmt.Sprintf("%d/%d embeddings successfully retrieved.", len(embeddings), len(body.Ids)),
			map[string]any{"embeddings": embeddings},
			nil,
		)
	}
}

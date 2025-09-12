package embeddings

import (
	"eigen_db/api/utils"
	"eigen_db/index"
	"eigen_db/types"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type deleteRequestBody struct {
	Ids []types.EmbId `json:"ids" binding:"required"`
}

func Delete(c *gin.Context) {
	var body deleteRequestBody
	if err := utils.ValidateBody(c, &body); err != nil {
		return
	}

	embeddingsDeleted := 0
	errors := make([]string, 0)
	for _, id := range body.Ids {
		if err := index.GetMemoryIndex().Delete(id); err != nil {
			errors = append(errors, fmt.Sprintf("embedding with ID %d was not deleted - %s", id, err.Error()))
		} else {
			embeddingsDeleted++
		}
	}

	if len(errors) != 0 {
		utils.SendResponse(
			c,
			http.StatusInternalServerError,
			fmt.Sprintf("%d/%d embeddings successfully deleted.", embeddingsDeleted, len(body.Ids)),
			nil,
			utils.CreateError("EMBEDDINGS_SKIPPED", errors),
		)
	} else {
		utils.SendResponse(
			c,
			200,
			fmt.Sprintf("%d/%d embeddings successfully deleted.", embeddingsDeleted, len(body.Ids)),
			nil,
			nil,
		)
	}
}

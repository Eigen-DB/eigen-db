package vector

import (
	"eigen_db/api/utils"
	"eigen_db/vector_io"
	"net/http"

	t "eigen_db/types"

	"github.com/gin-gonic/gin"
)

type insertRequestBody struct {
	Embedding t.Embedding `json:"embedding" binding:"required"`
}

func Insert(c *gin.Context) {
	var body insertRequestBody
	if err := utils.ValidateBody(c, &body); err != nil {
		return
	}

	v, err := vector_io.NewVector(body.Embedding)
	if err != nil {
		utils.SendResponse(
			c,
			http.StatusBadRequest,
			"The vector you provided is invalid.",
			nil,
			utils.CreateError("INVALID_VECTOR_PROVIDED", err.Error()),
		)
		return
	}

	if err := vector_io.InsertVector(v); err != nil { // causes nil pointer deference bug when empty body
		utils.SendResponse(
			c,
			http.StatusInternalServerError,
			"An error occured when inserting your vector.",
			nil,
			utils.CreateError("CANNOT_INSERT_VECTOR", err.Error()),
		)
		return
	}

	utils.SendResponse(
		c,
		http.StatusOK,
		"Vector successfully inserted.",
		nil,
		nil,
	)
}

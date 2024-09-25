package vector

import (
	"eigen_db/api/utils"
	"eigen_db/vector_io"
	"net/http"

	t "eigen_db/types"

	"github.com/gin-gonic/gin"
)

type insertRequestBody struct {
	Components []t.VectorComponent `json:"components" binding:"required"`
}

func Insert(vectorFactory vector_io.IVectorFactory) func(*gin.Context) {
	return func(c *gin.Context) {
		var body insertRequestBody
		if err := utils.ValidateBody(c, &body); err != nil {
			return
		}

		v, err := vectorFactory.NewVector(body.Components)
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

		if err := v.Insert(); err != nil { // causes nil pointer deference bug when empty body
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
}

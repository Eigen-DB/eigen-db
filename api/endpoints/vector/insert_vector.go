package vector

import (
	"eigen_db/vector_io"
	"net/http"

	t "eigen_db/types"

	"github.com/gin-gonic/gin"
)

type insertRequestBody struct {
	Components t.VectorComponents `json:"components" binding:"required"`
}

func Insert(vectorFactory vector_io.IVectorFactory) func(*gin.Context) {
	return func(c *gin.Context) {
		var body insertRequestBody
		if err := c.ShouldBindJSON(&body); err != nil {
			c.Status(http.StatusBadRequest)
			return
		}

		v, err := vectorFactory.NewVector(body.Components)
		if err != nil {
			c.String(http.StatusBadRequest, "vector provided had the wrong dimensionality")
			return
		}
		v.Insert() // causes nil pointer deference bug when empty body

		c.String(200, "Vector successfully inserted.")
	}
}

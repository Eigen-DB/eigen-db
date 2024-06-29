package vector

import (
	t "eigen_db/types"
	"eigen_db/vector_io"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type bulkInsertRequestBody struct {
	SetOfComponents []t.VectorComponents `json:"setOfComponents" binding:"required"`
}

func BulkInsert(vectorFactory vector_io.IVectorFactory) func(*gin.Context) {
	return func(c *gin.Context) {
		var body bulkInsertRequestBody
		if err := c.ShouldBindJSON(&body); err != nil {
			c.Status(http.StatusBadRequest)
			return
		}

		vectorsInserted := 0
		for i, components := range body.SetOfComponents {
			v, err := vectorFactory.NewVector(components)
			if err != nil {
				c.Error(fmt.Errorf("vector %d was skipped as it had the wrong dimensionality", i))
				continue
			}
			v.Insert()
			vectorsInserted++
		}

		response := fmt.Sprintf("%d/%d vectors successfully insertedx.", vectorsInserted, len(body.SetOfComponents))
		c.String(http.StatusOK, response) // later handle returning any vectors that were skipped to the client
	}
}

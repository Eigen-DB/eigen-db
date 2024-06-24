package vector_endpoints

import (
	t "eigen_db/types"
	"eigen_db/vector_io"
	"encoding/json"
	"io"

	"github.com/gin-gonic/gin"
)

type bulkInsertRequestBody struct {
	SetOfComponents []t.VectorComponents `json:"setOfComponents"`
}

func BulkInsertVector(c *gin.Context) {
	bodyBytes, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.Error(err)
	}

	body := &bulkInsertRequestBody{}
	err = json.Unmarshal(bodyBytes, body)
	if err != nil {
		c.Error(err)
	}

	var v *vector_io.Vector
	for _, c := range body.SetOfComponents {
		v = vector_io.NewVector(c)
		v.Insert()
	}

	c.String(200, "Vectors successfully bulk-inserted.")
}

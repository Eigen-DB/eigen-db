package api

import (
	vio "eigen_db/vector_io"
	"encoding/json"
	"io"

	"github.com/gin-gonic/gin"
)

type bulkInsertRequestBody struct {
	SetOfComponents []vio.VectorComponents `json:"setOfComponents"`
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

	var v *vio.Vector
	for _, c := range body.SetOfComponents {
		v = vio.NewVector(c)
		v.Insert()
	}

	c.String(200, "Vectors successfully bulk-inserted.")
}

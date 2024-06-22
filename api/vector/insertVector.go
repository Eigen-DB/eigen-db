package api

import (
	"eigen_db/vector_io"
	"encoding/json"
	"io"

	t "eigen_db/types"

	"github.com/gin-gonic/gin"
)

type insertRequestBody struct {
	Components t.VectorComponents `json:"components"`
}

func InsertVector(c *gin.Context) {
	bodyBytes, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.Error(err)
	}

	body := &insertRequestBody{}
	err = json.Unmarshal(bodyBytes, body)
	if err != nil {
		c.Error(err)
	}
	v := vector_io.NewVector(body.Components)
	v.Insert()

	c.String(200, "Vector successfully inserted.")
}

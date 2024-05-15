package api

import (
	"eigen_db/vector_io"
	"encoding/json"
	"io"

	"github.com/gin-gonic/gin"
)

type requestBody struct {
	Components []float64 `json:"components"`
}

func InsertVector(c *gin.Context) {
	bodyBytes, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.Error(err)
	}

	body := &requestBody{}
	err = json.Unmarshal(bodyBytes, body)
	if err != nil {
		c.Error(err)
	}
	v := vector_io.NewVector(body.Components)
	v.Insert()

	c.String(200, "Vector successfully inserted.")
}

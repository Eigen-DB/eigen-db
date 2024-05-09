package routes

import (
	"eigen_db/vectors"
	"encoding/json"

	"github.com/gin-gonic/gin"
)

type requestBody struct {
	components []float32
}

func InsertVector(c *gin.Context) {
	bodyBytes := make([]byte, 0)
	_, err := c.Request.Body.Read(bodyBytes)
	if err != nil {
		panic(err)
	}

	body := new(requestBody)
	err = json.Unmarshal(bodyBytes, body)
	if err != nil {
		panic(err)
	}
	v := vectors.NewVector(body.components)
	v.Insert()
}

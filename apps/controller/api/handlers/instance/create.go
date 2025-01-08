package instance

import (
	"controller/api/types"
	"controller/utils/k8s"
	"log"

	"github.com/gin-gonic/gin"
)

func CreateInstance(c *gin.Context) {
	var body types.CreateBody
	if err := c.ShouldBindJSON(&body); err != nil {
		log.Println(err.Error())
		c.String(400, err.Error())
		return
	}

	_, apiKey, err := k8s.SpawnInstance(body.CustomerId)
	if err != nil {
		log.Println(err.Error())
		c.String(400, err.Error())
		return
	}

	c.String(200, apiKey)
}

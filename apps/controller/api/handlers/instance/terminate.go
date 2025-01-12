package instance

import (
	"controller/api/types"
	"controller/utils/k8s"
	"log"

	"github.com/gin-gonic/gin"
)

func TerminateInstance(c *gin.Context) {
	var body types.TerminateBody
	if err := c.ShouldBindJSON(&body); err != nil {
		log.Println(err.Error())
		c.String(400, err.Error())
		return
	}

	_, err := k8s.TerminateInstance(body.CustomerId)
	if err != nil {
		log.Println(err.Error())
		c.String(400, err.Error())
		return
	}

	c.String(200, "Instance terminated.")
}

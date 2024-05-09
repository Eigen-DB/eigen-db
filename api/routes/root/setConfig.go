package routes

import "github.com/gin-gonic/gin"

type requestBody struct {
	similarityAlgorithm string
	dimensions          uint32
}

func SetConfig(c *gin.Context) {
	c.JSON(200, gin.H{
		"fu": "bar",
	})
}

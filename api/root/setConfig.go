package api

import (
	vio "eigen_db/vector_io"

	"github.com/gin-gonic/gin"
)

type requestBody struct {
	config vio.VectorSpaceConfig
}

func SetConfig(c *gin.Context) {
	c.JSON(200, gin.H{
		"fu": "bar",
	})
}

package api

import (
	"eigen_db/cfg"

	"github.com/gin-gonic/gin"
)

type requestBody struct {
	config cfg.Config
}

func SetConfig(c *gin.Context) {
	c.JSON(200, gin.H{
		"fu": "bar",
	})
}

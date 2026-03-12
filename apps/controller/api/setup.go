package api

import (
	"controller/api/handlers"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()
	v1 := r.Group("/api/v1")
	v1.POST("/start", handlers.Start)
	v1.POST("/stop", handlers.Stop)

	return r
}

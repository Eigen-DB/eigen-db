package main

import (
	"eigen_db/api/routes"
	"eigen_db/vectors"

	"github.com/gin-gonic/gin"
)

func setupRouter() *gin.Engine {
	r := gin.Default()

	r.GET("/ping", routes.Ping)
	r.POST("/set-config", routes.SetConfig)

	return r
}

func main() {
	vectors.InitializeVectorStorage()
	r := setupRouter()
	r.Run("127.0.0.1:8080")
}

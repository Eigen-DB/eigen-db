package api

import (
	"eigen_db/api/endpoints/root"
	"eigen_db/api/endpoints/vector"
	"eigen_db/vector_io"

	"github.com/gin-gonic/gin"
)

func setupRouter() *gin.Engine {
	r := gin.Default()
	vectors := r.Group("/vector")

	r.GET("/ping", root.Ping)
	r.POST("/update-config", root.UpdateConfig)
	vectors.PUT("/insert", vector.Insert(&vector_io.VectorFactory{}))
	vectors.PUT("/bulk-insert", vector.BulkInsert(&vector_io.VectorFactory{}))
	vectors.GET("/search", vector.Search(&vector_io.VectorSearcher{}))

	return r
}

func StartAPI(addr string) error {
	r := setupRouter()
	err := r.Run(addr)
	return err
}

package api

import (
	root_endpoints "eigen_db/api/root"
	vector_endpoints "eigen_db/api/vector"
	"eigen_db/vector_io"

	"github.com/gin-gonic/gin"
)

func setupRouter() *gin.Engine {
	r := gin.Default()
	vectors := r.Group("/vector")

	r.GET("/ping", root_endpoints.Ping)
	r.PUT("/update-config", root_endpoints.UpdateConfig)
	vectors.PUT("/insert", vector_endpoints.Insert(&vector_io.VectorFactory{}))
	vectors.PUT("/bulk-insert", vector_endpoints.BulkInsert(&vector_io.VectorFactory{}))
	vectors.GET("/search", vector_endpoints.Search)

	return r
}

func StartAPI(addr string) error {
	r := setupRouter()
	err := r.Run(addr)
	return err
}

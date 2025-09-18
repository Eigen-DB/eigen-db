package api

import (
	"eigen_db/api/middleware"
	v1_embeddings "eigen_db/api/v1/embeddings"
	v1_health_check "eigen_db/api/v1/health_check"
	v1_indexes "eigen_db/api/v1/indexes"

	"github.com/gin-gonic/gin"
)

// Setups up the API router
//
// Returns the router as a pointer to a Gin Engine instance.
func setupRouter() *gin.Engine {
	r := gin.Default()

	v1 := r.Group("/api/v1")

	// ***** START OF V1 ENDPOINTS *****
	embeddingsEndpoints := v1.Group("/embeddings", middleware.AuthMiddleware())
	indexesEndpoints := v1.Group("/indexes", middleware.AuthMiddleware())

	v1.GET("/health", v1_health_check.Health)
	v1.GET("/test-auth", middleware.AuthMiddleware(), v1_health_check.TestAuth)

	embeddingsEndpoints.PUT("/:index/insert", v1_embeddings.Insert)
	embeddingsEndpoints.PUT("/:index/upsert", v1_embeddings.Upsert)
	embeddingsEndpoints.DELETE("/:index/delete", v1_embeddings.Delete)
	embeddingsEndpoints.POST("/:index/retrieve", v1_embeddings.Retrieve)
	embeddingsEndpoints.POST("/:index/search", v1_embeddings.Search)

	indexesEndpoints.PUT("/:index/create", v1_indexes.Create)
	indexesEndpoints.DELETE("/:index/delete", v1_indexes.Delete)
	indexesEndpoints.GET("/:index/stats", v1_indexes.Stats)
	indexesEndpoints.GET("/list", v1_indexes.List)
	// ***** END OF V1 ENDPOINTS *****

	return r
}

// Starts the API server
//
// Returns an error if one occured.
func StartAPI(addr string) error {
	r := setupRouter()
	err := r.Run(addr)
	return err
}

package api

import (
	"eigen_db/api/endpoints/embeddings"
	"eigen_db/api/endpoints/health_check"
	"eigen_db/api/endpoints/indexes"
	"eigen_db/api/middleware"

	"github.com/gin-gonic/gin"
)

// Setups up the API router
//
// Returns the router as a pointer to a Gin Engine instance.
func setupRouter() *gin.Engine {
	r := gin.Default()

	embeddingsEndpoints := r.Group("/embeddings", middleware.AuthMiddleware())
	indexesEndpoints := r.Group("/indexes", middleware.AuthMiddleware())

	// updateConfigRoot := r.Group("/update-config", middleware.AuthMiddleware())

	// updatePersistence := updateConfigRoot.Group("/persistence")
	// updateApi := updateConfigRoot.Group("/api")
	// indexConfig := updateConfigRoot.Group("/hnsw-params")

	// health check endpoints
	r.GET("/health", health_check.Health)
	r.GET("/test-auth", middleware.AuthMiddleware(), health_check.TestAuth)

	// embedding management endpoints
	embeddingsEndpoints.PUT("/:index/insert", embeddings.Insert)
	embeddingsEndpoints.PUT("/:index/upsert", embeddings.Upsert)
	embeddingsEndpoints.DELETE("/:index/delete", embeddings.Delete)
	embeddingsEndpoints.POST("/:index/retrieve", embeddings.Retrieve)
	embeddingsEndpoints.POST("/:index/search", embeddings.Search)

	// index management endpoints
	indexesEndpoints.PUT("/:index/create", indexes.Create)
	indexesEndpoints.DELETE("/:index/delete", indexes.Delete)
	indexesEndpoints.GET("/:index/stats", indexes.Stats)
	indexesEndpoints.GET("/list", indexes.List)

	// config setter endpoints
	// updatePersistence.POST("/time-interval", persistence.UpdateTimeInterval)
	// updateApi.POST("/port", api.UpdatePort)
	// updateApi.POST("/address", api.UpdateAddress)
	// indexConfig.POST("/similarity-metric", index_config.UpdateSimilarityMetric)

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

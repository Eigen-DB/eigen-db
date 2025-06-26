package api

import (
	"eigen_db/api/endpoints/embeddings"
	"eigen_db/api/endpoints/health_check"
	"eigen_db/api/endpoints/update_config/api"
	"eigen_db/api/endpoints/update_config/index_config"
	"eigen_db/api/endpoints/update_config/persistence"
	"eigen_db/api/middleware"

	"github.com/gin-gonic/gin"
)

// Setups up the API router
//
// Returns the router at a pointer to a Gin Engine instance.
func setupRouter() *gin.Engine {
	r := gin.Default()

	vectors := r.Group("/embeddings", middleware.AuthMiddleware())
	updateConfigRoot := r.Group("/update-config", middleware.AuthMiddleware())

	updatePersistence := updateConfigRoot.Group("/persistence")
	updateApi := updateConfigRoot.Group("/api")
	indexConfig := updateConfigRoot.Group("/hnsw-params")

	// health check endpoints
	r.GET("/health", health_check.Health)
	r.GET("/test-auth", middleware.AuthMiddleware(), health_check.TestAuth)

	// vector operation endpoints
	vectors.PUT("/insert", embeddings.Insert)
	vectors.PUT("/upsert", embeddings.Upsert)
	vectors.GET("/retrieve", embeddings.Retrieve)
	vectors.GET("/search", embeddings.Search)
	// config setter endpoints
	updatePersistence.POST("/time-interval", persistence.UpdateTimeInterval)
	updateApi.POST("/port", api.UpdatePort)
	updateApi.POST("/address", api.UpdateAddress)
	indexConfig.POST("/similarity-metric", index_config.UpdateSimilarityMetric)

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

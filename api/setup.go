package api

import (
	"eigen_db/api/endpoints/health_check"
	"eigen_db/api/endpoints/update_config/api"
	"eigen_db/api/endpoints/update_config/hnsw_params"
	"eigen_db/api/endpoints/update_config/persistence"
	"eigen_db/api/endpoints/vector"
	"eigen_db/api/middleware"

	"github.com/gin-gonic/gin"
)

// Setups up the API router
//
// Returns the router at a pointer to a Gin Engine instance.
func setupRouter() *gin.Engine {
	r := gin.Default()

	vectors := r.Group("/vector", middleware.AuthMiddleware())
	updateConfigRoot := r.Group("/update-config", middleware.AuthMiddleware())

	updatePersistence := updateConfigRoot.Group("/persistence")
	updateApi := updateConfigRoot.Group("/api")
	updateHnswParams := updateConfigRoot.Group("/hnsw-params")

	// health check endpoints
	r.GET("/health", health_check.Health)
	r.POST("/test-auth", middleware.AuthMiddleware(), health_check.TestAuth)

	// vector operation endpoints
	vectors.PUT("/insert", vector.Insert)
	vectors.PUT("/bulk-insert", vector.BulkInsert)
	vectors.GET("/search", vector.Search)

	// config setter endpoints
	updatePersistence.POST("/time-interval", persistence.UpdateTimeInterval)
	updateApi.POST("/port", api.UpdatePort)
	updateApi.POST("/address", api.UpdateAddress)
	updateHnswParams.POST("/similarity-metric", hnsw_params.UpdateSimilarityMetric)
	updateHnswParams.POST("/vector-space-size", hnsw_params.UpdateSpaceSize)
	updateHnswParams.POST("/m", hnsw_params.UpdateM)
	updateHnswParams.POST("/ef-construction", hnsw_params.UpdateEfConstruction)

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

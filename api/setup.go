package api

import (
	"context"
	"eigen_db/api/endpoints/health_check"
	"eigen_db/api/endpoints/update_config/api"
	"eigen_db/api/endpoints/update_config/hnsw_params"
	"eigen_db/api/endpoints/update_config/persistence"
	"eigen_db/api/endpoints/vector"
	"eigen_db/api/middleware"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func setupRouter(ctx context.Context, redisClient *redis.Client) *gin.Engine {
	r := gin.Default()

	vectors := r.Group("/vector", middleware.AuthMiddleware(ctx, redisClient))
	updateConfigRoot := r.Group("/update-config", middleware.AuthMiddleware(ctx, redisClient))

	updatePersistence := updateConfigRoot.Group("/persistence")
	updateApi := updateConfigRoot.Group("/api")
	updateHnswParams := updateConfigRoot.Group("/hnsw-params")

	// health check endpoint
	r.GET("/health", health_check.Health)

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

func StartAPI(ctx context.Context, addr string, redisClient *redis.Client) error {
	r := setupRouter(ctx, redisClient)
	err := r.Run(addr)
	return err
}

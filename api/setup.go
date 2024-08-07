package api

import (
	"eigen_db/api/endpoints/health_check"
	"eigen_db/api/endpoints/update_config/api"
	"eigen_db/api/endpoints/update_config/hnsw_params"
	"eigen_db/api/endpoints/update_config/persistence"
	"eigen_db/api/endpoints/vector"
	"eigen_db/cfg"
	"eigen_db/vector_io"

	"github.com/gin-gonic/gin"
)

func setupRouter() *gin.Engine {
	config := (&cfg.ConfigFactory{}).GetConfig()

	r := gin.Default()

	vectors := r.Group("/vector")

	updateConfigRoot := r.Group("/update-config")
	updatePersistence := updateConfigRoot.Group("/persistence")
	updateApi := updateConfigRoot.Group("/api")
	updateHnswParams := updateConfigRoot.Group("/hnsw-params")

	// health check endpoint
	r.GET("/health", health_check.Health)

	// vector operation endpoints
	vectors.PUT("/insert", vector.Insert(&vector_io.VectorFactory{}))
	vectors.PUT("/bulk-insert", vector.BulkInsert(&vector_io.VectorFactory{}))
	vectors.GET("/search", vector.Search(&vector_io.VectorSearcher{}))

	// config setter endpoints
	updatePersistence.POST("/time-interval", persistence.UpdateTimeInterval(config))

	updateApi.POST("/port", api.UpdatePort(config))
	updateApi.POST("/address", api.UpdateAddress(config))

	updateHnswParams.POST("/dimensions", hnsw_params.UpdateDimensions(config))
	updateHnswParams.POST("/similarity-metric", hnsw_params.UpdateSimilarityMetric(config))
	updateHnswParams.POST("/vector-space-size", hnsw_params.UpdateSpaceSize(config))
	updateHnswParams.POST("/m", hnsw_params.UpdateM(config))
	updateHnswParams.POST("/ef-construction", hnsw_params.UpdateEfConstruction(config))

	return r
}

func StartAPI(addr string) error {
	r := setupRouter()
	err := r.Run(addr)
	return err
}

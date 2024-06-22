package main

import (
	root "eigen_db/api/root"
	vector "eigen_db/api/vector"
	"eigen_db/cfg"
	"eigen_db/vector_io"
	"fmt"

	"github.com/gin-gonic/gin"
)

func setupAPIRouter() *gin.Engine {
	r := gin.Default()
	vectors := r.Group("/vector")

	r.GET("/ping", root.Ping)
	r.POST("/set-config", root.SetConfig)
	vectors.POST("/insert", vector.InsertVector)
	vectors.POST("/bulk-insert", vector.BulkInsertVector)
	vectors.POST("/search", vector.Search)

	return r
}

func startAPI(addr string) error {
	r := setupAPIRouter()
	err := r.Run(addr)
	return err
}

func setupDB(config *cfg.Config) {
	params := config.HNSWParams
	vector_io.InstantiateVectorStore(
		params.Dimensions,
		params.SimilarityMetric,
		params.SpaceSize,
		params.M,
		params.EfConstruction,
	)
}

func main() {
	cfg.UpdateConfig() // parses config.yml into memory
	config := cfg.GetConfig()

	setupDB(config)

	if err := vector_io.StartPersistenceLoop(config); err != nil {
		panic(err)
	}

	if err := startAPI(fmt.Sprintf("%s:%d", config.API.Address, config.API.Port)); err != nil {
		panic(err)
	}
}

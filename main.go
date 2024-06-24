package main

import (
	root_endpoints "eigen_db/api/root"
	vector_endpoints "eigen_db/api/vector"
	"eigen_db/cfg"
	"eigen_db/vector_io"
	"fmt"

	"github.com/gin-gonic/gin"
)

func setupAPIRouter() *gin.Engine {
	r := gin.Default()
	vectors := r.Group("/vector")

	r.GET("/ping", root_endpoints.Ping)
	r.POST("/set-config", root_endpoints.SetConfig)
	vectors.POST("/insert", vector_endpoints.InsertVector)
	vectors.POST("/bulk-insert", vector_endpoints.BulkInsertVector)
	vectors.POST("/search", vector_endpoints.Search)

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
	cfg.NewConfig()           // creates a empty Config struct in memory
	config := cfg.GetConfig() // get pointer to Config in memory
	config.LoadConfig()       // load config from config.yml into the Config struct in memory

	setupDB(config)

	if err := vector_io.StartPersistenceLoop(config); err != nil {
		panic(err)
	}

	if err := startAPI(fmt.Sprintf("%s:%d", config.API.Address, config.API.Port)); err != nil {
		panic(err)
	}
}

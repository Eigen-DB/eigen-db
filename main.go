package main

import (
	root "eigen_db/api/root"
	vector "eigen_db/api/vector"
	"eigen_db/vector_io"

	"github.com/gin-gonic/gin"
)

const enablePersistence = true // turn this on to enable the 5 second persistence loop

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

func setupDB() {
	vector_io.InstantiateVectorStore(enablePersistence, 2, vector_io.EUCLIDEAN, 10000, 32, 400) // use real params
}

func main() {
	setupDB()

	if enablePersistence {
		if err := vector_io.StartPersistenceLoop(); err != nil {
			panic(err)
		}
	}

	if err := startAPI("0.0.0.0:8080"); err != nil {
		panic(err)
	}
}

package main

import (
	root "eigen_db/api/routes/root"
	vector "eigen_db/api/routes/vector"
	c "eigen_db/constants"
	"os"

	"github.com/gin-gonic/gin"
)

func setupRouter() *gin.Engine {
	r := gin.Default()
	vectors := r.Group("/vector")

	r.GET("/ping", root.Ping)
	r.POST("/set-config", root.SetConfig)
	vectors.POST("/insert", vector.InsertVector)

	return r
}

func startAPI(addr string) {
	r := setupRouter()
	r.Run(addr)
}

func setupDB() {
	err := os.Mkdir(c.PERSISTENT_VOLUME_PATH, 0700) // rwx------
	if err != nil {
		if !os.IsExist(err) {
			panic(err)
		}
	}

	if _, err := os.Stat(c.DATABASE_PATH); os.IsNotExist(err) { // if database does not exist
		file, err := os.Create(c.DATABASE_PATH)
		// initialize the JSON...
		if err != nil {
			panic(err)
		}
		defer file.Close()
	}

	if _, err := os.Stat(c.CONFIG_PATH); os.IsNotExist(err) { // if config does not exist
		file, err := os.Create(c.DATABASE_PATH)
		// initialize the JSON...
		if err != nil {
			panic(err)
		}
		defer file.Close()
	}
}

func main() {
	setupDB()
	startAPI("127.0.0.1:8080")
}

package main

import (
	"context"
	"eigen_db/api"
	"eigen_db/cfg"
	"eigen_db/redis_utils"
	"eigen_db/vector_io"
	"flag"
	"fmt"
	"os"
)

func main() {
	var apiKey string
	var redisHost string
	var redisPort string
	var redisPass string

	flag.StringVar(&apiKey, "api-key", "", "EigenDB API key")
	flag.StringVar(&redisHost, "redis-host", "127.0.0.1", "Redis server host IP (default: 127.0.0.1)")
	flag.StringVar(&redisPort, "redis-port", "6379", "Redis server host port (default: 6379)")
	flag.StringVar(&redisPass, "redis-pass", "", "Redis server password (default: \"\")")
	flag.Parse()

	// setting up the in-memory config
	if err := cfg.SetupConfig(); err != nil {
		fmt.Println("There was an error with setting up the config")
		panic(err)
	}
	config := cfg.GetConfig()

	// checking if EigenDB is running in TEST_MODE
	if os.Getenv("TEST_MODE") == "1" {
		fmt.Println("*** EigenDB running in TEST_MODE, making the API key = \"test\". If this was not intentional, please run EigenDB in standard mode. ***")
		apiKey = "test"
		if err := config.SetHNSWParamsDimensions(2); err != nil { // setting dimensions to 2 for the tests
			fmt.Println("An error occured when setting the dimensions to 2.")
			panic(err)
		}
	}

	// setting up the in-memory vector store
	if err := vector_io.InstantiateVectorStore(
		config.GetHNSWParamsDimensions(),
		config.GetHNSWParamsSimilarityMetric(),
		config.GetHNSWParamsSpaceSize(),
		config.GetHNSWParamsM(),
		config.GetHNSWParamsEfConstruction(),
	); err != nil {
		panic(err)
	}

	//setting up the Redis connection
	ctx := context.Background()
	os.Setenv("REDIS_HOST", redisHost)
	os.Setenv("REDIS_PORT", redisPort)
	os.Setenv("REDIS_PASS", redisPass)
	redisClient, err := redis_utils.GetConnection(ctx)
	if err != nil {
		panic(err)
	}

	// setting up the API key
	apiKey, err = redis_utils.SetupAPIKey(ctx, redisClient, apiKey)
	if err != nil {
		panic(err)
	}
	fmt.Printf("API KEY: %s\n", apiKey)

	// starting the persistence loop
	if err := vector_io.StartPersistenceLoop(config); err != nil {
		panic(err)
	}

	// setting up the REST API
	if err := api.StartAPI(ctx, fmt.Sprintf("%s:%d", config.GetAPIAddress(), config.GetAPIPort()), redisClient); err != nil {
		panic(err)
	} else {
		fmt.Println(apiKey)
	}
}

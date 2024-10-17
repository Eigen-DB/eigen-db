package main

import (
	"eigen_db/api"
	"eigen_db/auth"
	"eigen_db/cfg"
	"eigen_db/constants"
	"eigen_db/metrics"
	"eigen_db/vector_io"
	"flag"
	"fmt"
	"os"
)

func main() {
	// initialize the metrics
	metrics.Init()

	// parsing cmd line args
	var apiKey string
	var regenApiKey bool

	flag.StringVar(&apiKey, "api-key", "", "EigenDB API key")
	flag.BoolVar(&regenApiKey, "regen-api-key", false, "Regenerate the API key")
	flag.Parse()

	// setting up the in-memory config
	if err := cfg.SetupConfig(constants.CONFIG_PATH); err != nil {
		fmt.Println("There was an error with setting up the config")
		panic(err)
	}
	config := cfg.GetConfig()

	// checking if EigenDB is running in TEST_MODE
	if os.Getenv("TEST_MODE") == "1" {
		fmt.Println("*** EigenDB running in TEST_MODE, making the API key = \"test\". If this was not intentional, please run EigenDB in standard mode. ***")
		apiKey = "test"
		if err := config.SetDimensions(2); err != nil { // setting dimensions to 2 for the tests
			fmt.Println("An error occured when setting the dimensions to 2.")
			panic(err)
		}
	}

	// setting up the in-memory vector store
	if err := vector_io.InstantiateVectorStore(
		config.GetDimensions(),
		config.GetSimilarityMetric(),
		config.GetSpaceSize(),
		config.GetM(),
		config.GetEfConstruction(),
	); err != nil {
		panic(err)
	}

	// setting up the API key
	apiKey, err := auth.SetupAPIKey(apiKey, regenApiKey, constants.API_KEY_FILE_PATH)
	if err != nil {
		panic(err)
	}
	fmt.Printf("API KEY: %s\n", apiKey)

	// starting the persistence loop
	if err := vector_io.StartPersistenceLoop(config); err != nil {
		panic(err)
	}

	// setting up the REST API
	if err := api.StartAPI(fmt.Sprintf("%s:%d", config.GetAPIAddress(), config.GetAPIPort())); err != nil {
		panic(err)
	} else {
		fmt.Println(apiKey)
	}
}

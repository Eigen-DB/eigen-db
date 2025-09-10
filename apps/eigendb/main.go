package main

import (
	"eigen_db/api"
	"eigen_db/auth"
	"eigen_db/cfg"
	"eigen_db/constants"
	"eigen_db/metrics"
	"eigen_db/types"
	"eigen_db/vector_io"
	"flag"
	"fmt"
	"os"
	"time"
)

func main() {
	e2e_test_mode := os.Getenv("E2E_TEST_MODE") == "1"

	// initialize the metrics
	metrics.Init()

	// parsing cmd line args
	var apiKey string
	var regenApiKey bool
	var persistenceTimeInterval time.Duration
	var apiPort int
	var apiAddress string
	var indexDimensions int
	var indexSimilarityMetric string

	flag.StringVar(&apiKey, "api-key", "", "EigenDB API key")
	flag.BoolVar(&regenApiKey, "regen-api-key", false, "Regenerate the API key")

	flag.DurationVar(&persistenceTimeInterval, "persistence-time-interval", time.Duration(0), "How often should data be persisted to disk (secs)")
	flag.IntVar(&apiPort, "api-port", 0, "API port")
	flag.StringVar(&apiAddress, "api-address", "", "API address")
	flag.IntVar(&indexDimensions, "dimensions", 0, "Dimensions")
	flag.StringVar(&indexSimilarityMetric, "similarity-metric", "", "Similarity metric")
	flag.Parse()

	// checking if EigenDB is running in E2E_TEST_MODE
	if e2e_test_mode {
		fmt.Println("*** EigenDB running in E2E_TEST_MODE, if this was not intentional, please run EigenDB in standard mode. ***\nSetting the API key = \"test\"")
		apiKey = "test"
	}

	// setting up the in-memory config
	if err := cfg.SetupConfig(constants.CONFIG_PATH); err != nil {
		fmt.Println("There was an error with setting up the config")
		panic(err)
	}
	config := cfg.GetConfig() // getting the in-memory config

	// create a map of config setters
	configSetters := map[bool]func() error{
		persistenceTimeInterval != time.Duration(0): func() error { return config.SetPersistenceTimeInterval(persistenceTimeInterval) },
		apiPort != 0:         func() error { return config.SetAPIPort(apiPort) },
		apiAddress != "":     func() error { return config.SetAPIAddress(apiAddress) },
		indexDimensions != 0: func() error { return config.SetDimensions(indexDimensions) },
		indexSimilarityMetric != "": func() error {
			metric := types.SimMetric(indexSimilarityMetric)
			if err := metric.Validate(); err != nil {
				return err
			}
			return config.SetSimilarityMetric(metric)
		},
	}

	for condition, setter := range configSetters {
		if condition {
			fmt.Println("Overriding config value")
			if err := setter(); err != nil {
				panic(err)
			}
		}
	}

	if err := config.WriteToDisk(constants.CONFIG_PATH); err != nil { // persisted new config to disk
		panic(err)
	}

	// setting up the in-memory vector store
	if err := vector_io.MemoryIndexInit(
		config.GetDimensions(),
		config.GetSimilarityMetric(),
	); err != nil {
		panic(err)
	}

	// setting up the API key
	apiKey, err := auth.SetupAPIKey(apiKey, regenApiKey, constants.API_KEY_FILE_PATH)
	if err != nil {
		panic(err)
	}
	fmt.Printf("IMPORTANT: API key has been generated and saved to %s\n", constants.API_KEY_FILE_PATH)
	fmt.Printf("API KEY: %s\n", apiKey)

	// starting the persistence loop
	if err := vector_io.GetMemoryIndex().StartPersistenceLoop(config); err != nil {
		panic(err)
	}

	// setting up the REST API
	if err := api.StartAPI(fmt.Sprintf("%s:%d", config.GetAPIAddress(), config.GetAPIPort())); err != nil {
		panic(err)
	} else {
		fmt.Println(apiKey)
	}
}

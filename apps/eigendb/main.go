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
	"strconv"
	"time"
)

func main() {
	// initialize the metrics
	metrics.Init()

	// parsing cmd line args
	var apiKey string
	var regenApiKey bool
	var persistenceTimeInterval time.Duration
	var apiPort int
	var apiAddress string
	var hnswDimensions int
	var hnswSimilarityMetric string
	var hnswVectorSpaceSize string // take argument as string and cast as uint32
	var hnswM int
	var hnswEfConstruction int

	flag.StringVar(&apiKey, "api-key", "", "EigenDB API key")
	flag.BoolVar(&regenApiKey, "regen-api-key", false, "Regenerate the API key")

	flag.DurationVar(&persistenceTimeInterval, "persistence-time-interval", time.Duration(0), "How often should data be persisted to disk (secs)")
	flag.IntVar(&apiPort, "api-port", 0, "API port")
	flag.StringVar(&apiAddress, "api-address", "", "API address")
	flag.IntVar(&hnswDimensions, "dimensions", 0, "Dimensions")
	flag.StringVar(&hnswSimilarityMetric, "similarity-metric", "", "Similarity metric")
	flag.StringVar(&hnswVectorSpaceSize, "vector-space-size", "", "Vector space size")
	flag.IntVar(&hnswM, "m", 0, "m parameter")
	flag.IntVar(&hnswEfConstruction, "efConstruction", 0, "efConstruction parameter")
	flag.Parse()

	// checking if EigenDB is running in E2E_TEST_MODE
	if os.Getenv("E2E_TEST_MODE") == "1" {
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
		apiPort != 0:        func() error { return config.SetAPIPort(apiPort) },
		apiAddress != "":    func() error { return config.SetAPIAddress(apiAddress) },
		hnswDimensions != 0: func() error { return config.SetDimensions(hnswDimensions) },
		hnswSimilarityMetric != "": func() error {
			metric := types.SimMetric(hnswSimilarityMetric)
			if err := metric.Validate(); err != nil {
				return err
			}
			return config.SetSimilarityMetric(metric)
		},
		hnswVectorSpaceSize != "": func() error {
			spaceSize, err := strconv.ParseUint(hnswVectorSpaceSize, 10, 32)
			if err != nil {
				return err
			}
			return config.SetSpaceSize(uint32(spaceSize))
		},
		hnswM != 0:              func() error { return config.SetM(hnswM) },
		hnswEfConstruction != 0: func() error { return config.SetEfConstruction(hnswEfConstruction) },
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

package main

import (
	"eigen_db/api"
	"eigen_db/auth"
	"eigen_db/cfg"
	"eigen_db/constants"
	"eigen_db/index_mgr"
	"eigen_db/metrics"
	"flag"
	"fmt"
	"os"
	"sync"
	"time"
)

func main() {
	e2eTestMode := os.Getenv("E2E_TEST_MODE") == "1"

	// initialize the metrics
	metrics.Init()

	// parsing cmd line args
	var apiKey string
	var regenApiKey bool
	var persistenceTimeInterval time.Duration
	var apiPort int
	var apiAddress string

	flag.StringVar(&apiKey, "api-key", "", "EigenDB API key")
	flag.BoolVar(&regenApiKey, "regen-api-key", false, "Regenerate the API key")

	flag.DurationVar(&persistenceTimeInterval, "persistence-time-interval", time.Duration(0), "How often should data be persisted to disk (secs)")
	flag.IntVar(&apiPort, "api-port", 0, "API port")
	flag.StringVar(&apiAddress, "api-address", "", "API address")
	flag.Parse()

	// checking if EigenDB is running in E2E_TEST_MODE
	if e2eTestMode {
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
		apiPort != 0:     func() error { return config.SetAPIPort(apiPort) },
		apiAddress != "": func() error { return config.SetAPIAddress(apiAddress) },
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

	// setting up the API key
	apiKey, err := auth.SetupAPIKey(apiKey, regenApiKey, constants.API_KEY_FILE_PATH)
	if err != nil {
		panic(err)
	}
	fmt.Printf("IMPORTANT: API key has been generated and saved to %s\n", constants.API_KEY_FILE_PATH)
	fmt.Printf("API KEY: %s\n", apiKey)

	wg := new(sync.WaitGroup)
	wg.Add(1)
	// initialize the index manager in memory
	if err := index_mgr.IndexMgrInit(wg, e2eTestMode); err != nil {
		panic(err)
	}
	// load any persisted indexes from the disk into memory
	mgr := index_mgr.GetIndexMgr()
	if err := mgr.LoadIndexes(wg, e2eTestMode); err != nil {
		panic(err)
	}

	// setting up the REST API
	if err := api.StartAPI(fmt.Sprintf("%s:%d", config.GetAPIAddress(), config.GetAPIPort())); err != nil {
		panic(err)
	} else {
		fmt.Println(apiKey)
	}
}

package main

import (
	"context"
	"eigen_db/api"
	"eigen_db/cfg"
	"eigen_db/constants"
	"eigen_db/redis_utils"
	"eigen_db/vector_io"
	"flag"
	"fmt"
	"os"
)

func displayAsciiArt() { // because it looks cool
	fmt.Println(`                                                                                                                                                  
	@@@@@@,                                                           
	@@@@@@@@@@@@&                                                     
	        @@@@@@@@                                                  
	            @@@@@@,                                               
	               @@@@@                                              
	                 @@@@@                                            
	                  @@@@@                                           
	                   .@@@@@                                         
	                     @@@@@                                        
	                      /@@@@#                                      
	                        @@@@@                                     
	                         %@@@@*                                   
	                  @@@@@    @@@@@                                  
	                @@@@@@      @@@@@.                                
	    @@&      @@@@@@.          @@@@@                               
	   ,@@@@   @@@@@@              @@@@@                              
	   @@@@.@@@@@@.                  @@@@@                            
	  @@@@@@@@@@                      ,@@@@@@                         
	  @@@@@@@@@@@@@@.                    @@@@@@@@,                    
	 @@@@@@@@@@@@@@                         *@@@@@@@@@@@@             
	%@@@@                                          &@@@@@                                                               
	`)
}

func main() {
	var apiKey string
	var redisHost string
	var redisPort string
	var redisPass string
	var noBanner bool

	flag.StringVar(&apiKey, "api-key", "", "EigenDB API key")
	flag.StringVar(&redisHost, "redis-host", "127.0.0.1", "Redis server host IP (default: 127.0.0.1)")
	flag.StringVar(&redisPort, "redis-port", "6379", "Redis server host port (default: 6379)")
	flag.StringVar(&redisPass, "redis-pass", "", "Redis server password (default: \"\")")
	flag.BoolVar(&noBanner, "no-banner", false, "Remove the ASCII banner at startup (default: false)")
	flag.Parse()

	if !noBanner {
		displayAsciiArt()
	}

	cfg.InstantiateConfig()                                          // creates a empty Config struct in memory
	config := cfg.GetConfig()                                        // get pointer to Config in memory
	if err := config.LoadConfig(constants.CONFIG_PATH); err != nil { // load config from config.yml into the Config struct in memory
		fmt.Println("There was an error loading the config into memory.")
		panic(err)
	}

	if os.Getenv("TEST_MODE") == "1" {
		fmt.Println("*** EigenDB running in TEST MODE, making the API key = \"test\". If this was not intentional, please run EigenDB in standard mode. ***")
		apiKey = "test"
		if err := config.SetHNSWParamsDimensions(2); err != nil { // setting dimensions to 2 for the tests
			fmt.Println("An error occured when setting the dimensions to 2.")
			panic(err)
		}
	}

	if err := vector_io.InstantiateVectorStore(
		config.GetHNSWParamsDimensions(),
		config.GetHNSWParamsSimilarityMetric(),
		config.GetHNSWParamsSpaceSize(),
		config.GetHNSWParamsM(),
		config.GetHNSWParamsEfConstruction(),
	); err != nil {
		panic(err)
	}

	ctx := context.Background()
	os.Setenv("REDIS_HOST", redisHost)
	os.Setenv("REDIS_PORT", redisPort)
	os.Setenv("REDIS_PASS", redisPass)

	redisClient, err := redis_utils.GetConnection(ctx)
	if err != nil {
		panic(err)
	}

	apiKey, err = redis_utils.SetupAPIKey(ctx, redisClient, apiKey)
	if err != nil {
		panic(err)
	}
	fmt.Printf("API KEY: %s\n", apiKey)

	if err := vector_io.StartPersistenceLoop(config); err != nil {
		panic(err)
	}

	if err := api.StartAPI(ctx, fmt.Sprintf("%s:%d", config.GetAPIAddress(), config.GetAPIPort()), redisClient); err != nil {
		panic(err)
	} else {
		fmt.Println(apiKey)
	}
}

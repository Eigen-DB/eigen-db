package main

/*
#cgo LDFLAGS: -L./lib -lhnsw
*/
import "C"

import (
	"eigen_db/api"
	"eigen_db/cfg"
	"eigen_db/constants"
	"eigen_db/vector_io"
	"fmt"
)

func main() {
	cfg.NewConfig()                          // creates a empty Config struct in memory
	config := cfg.GetConfig()                // get pointer to Config in memory
	config.LoadConfig(constants.CONFIG_PATH) // load config from config.yml into the Config struct in memory

	vector_io.SetupDB(config)

	if err := vector_io.StartPersistenceLoop(config); err != nil {
		panic(err)
	}

	if err := api.StartAPI(fmt.Sprintf("%s:%d", config.API.Address, config.API.Port)); err != nil {
		panic(err)
	}
}

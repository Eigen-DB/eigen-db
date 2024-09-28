package cfg

import (
	"eigen_db/constants"
)

func SetupConfig() error {
	instantiateConfig()                                                  // creates a empty Config struct in memory
	config := GetConfig()                                                // get pointer to Config in memory
	if err := config.populateConfig(constants.CONFIG_PATH); err != nil { // populate config in memory with values from config.yml
		return err
	}
	return nil
}

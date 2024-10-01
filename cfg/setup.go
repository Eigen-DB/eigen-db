package cfg

func SetupConfig(configPath string) error {
	instantiateConfig()                                       // creates a empty Config struct in memory
	config := GetConfig()                                     // get pointer to Config in memory
	if err := config.populateConfig(configPath); err != nil { // populate config in memory with values from config.yml
		return err
	}
	return nil
}

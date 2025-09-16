package cfg

import (
	"eigen_db/constants"
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/manifoldco/promptui"
)

// Wrapper function for all the steps to properly set up the in-memory config at program start.
//
// Returns an error if one occured.
func SetupConfig(configPath string) error {
	instantiateConfig()                               // creates a empty Config struct in memory
	config := GetConfig()                             // get pointer to Config in memory
	if os.Getenv("EIGENDB_INTERACTIVE_MENU") == "1" { // if the user wants to use the interactive menu, start it (does not work very well with docker compose, only useful for local development)
		return startConfigMenu()
	}
	if _, err := os.Stat(configPath); errors.Is(err, os.ErrNotExist) { // if config file does not exist
		return errors.New("config file does not exist, please create one at " + configPath + " or run EigenDB with EIGENDB_INTERACTIVE_MENU=1 to create one using the interactive menu")
	}
	if err := config.populateConfig(configPath); err != nil { // populate config in memory with values from config.yml
		return err
	}
	return nil
}

func startConfigMenu() error {
	fmt.Println("Please select your configuration values.")
	config := GetConfig()

	// setting persistence time interval
	result, err := (&promptui.Prompt{
		Label: "How often should data be persisted to disk (secs) (>= 1s)",
		Validate: func(input string) error {
			val, err := strconv.ParseFloat(input, 32)
			if err != nil {
				return errors.New("persistence time interval must be a valid 32-bit float")
			}
			if val < 1 {
				return errors.New("persistence time interval must be >= 1s")
			}
			return nil
		},
		Default: "3",
	}).Run()
	if err != nil {
		return err
	}
	interval, _ := strconv.ParseFloat(result, 32) // ignoring error as the input is already validated as a valid float32 when received
	if err := config.SetPersistenceTimeInterval(time.Duration(interval * float64(time.Second))); err != nil {
		return err
	}

	// setting api port
	result, err = (&promptui.Prompt{
		Label: "API port",
		Validate: func(input string) error {
			val, err := validateInt32(input)
			if err != nil {
				return err
			}
			if val < 0 {
				return errors.New("value must be a valid port number")
			}
			return nil
		},
		Default: "8080",
	}).Run()
	if err != nil {
		return err
	}
	port, _ := strconv.ParseInt(result, 10, 32)
	if err := config.SetAPIPort(int(port)); err != nil {
		return err
	}

	// setting api address
	result, err = (&promptui.Prompt{
		Label:   "API address",
		Default: "0.0.0.0",
	}).Run()
	if err != nil {
		return err
	}
	if err := config.SetAPIAddress(result); err != nil {
		return err
	}

	return config.WriteToDisk(constants.CONFIG_PATH)
}

func validateInt32(input string) (int, error) {
	val, err := strconv.ParseInt(input, 10, 32)
	if err != nil {
		return 0, errors.New("value must be a valid integer")
	}
	return int(val), nil
}

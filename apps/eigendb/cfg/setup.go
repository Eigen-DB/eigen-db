package cfg

import (
	"eigen_db/constants"
	"eigen_db/types"
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
	instantiateConfig()   // creates a empty Config struct in memory
	config := GetConfig() // get pointer to Config in memory
	if os.Getenv("E2E_TEST_MODE") == "1" {
		fmt.Println("Making E2E test config")
		_ = config.SetPersistenceTimeInterval(3 * time.Second)
		_ = config.SetAPIPort(8080)
		_ = config.SetAPIAddress("0.0.0.0") // wouldn't 127.0.0.1 be better ?
		_ = config.SetDimensions(2)
		_ = config.SetSimilarityMetric(types.MetricL2)
		if err := config.WriteToDisk(constants.CONFIG_PATH); err != nil {
			return err
		}
		return nil
	}

	if _, err := os.Stat(configPath); errors.Is(err, os.ErrNotExist) { // if config file does not exist -> choose your config values
		fmt.Println("No existing config file found. Please select your configuration values.")
		return startConfigMenu()
	}
	if err := config.populateConfig(configPath); err != nil { // populate config in memory with values from config.yml
		return err
	}
	return nil
}

func startConfigMenu() error {
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

	// setting dimensions
	result, err = (&promptui.Prompt{
		Label: "Dimensions (>= 2)",
		Validate: func(input string) error {
			val, err := validateInt32(input)
			if err != nil {
				return err
			}
			if val < 2 {
				return errors.New("dimensions must be >= 2")
			}
			return nil
		},
	}).Run()
	if err != nil {
		return err
	}
	dim, _ := strconv.ParseInt(result, 10, 32)
	if err := config.SetDimensions(int(dim)); err != nil {
		return err
	}

	// setting similarity metric
	_, result, err = (&promptui.Select{
		Label: "Similarity metric",
		Items: []string{
			types.MetricCosine.String(),
			types.MetricInnerProduct.String(),
			types.MetricL2.String(),
		},
	}).Run()
	if err != nil {
		return err
	}
	if err := config.SetSimilarityMetric(types.SimMetric(result)); err != nil {
		return err
	}

	return config.WriteToDisk(constants.CONFIG_PATH) // persisting config values to disk
}

func validateInt32(input string) (int, error) {
	val, err := strconv.ParseInt(input, 10, 32)
	if err != nil {
		return 0, errors.New("value must be a valid integer")
	}
	return int(val), nil
}

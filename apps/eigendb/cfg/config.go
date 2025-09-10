package cfg

import (
	"errors"
	"os"
	"time"

	"eigen_db/constants"
	"eigen_db/types"
	t "eigen_db/types"

	"gopkg.in/yaml.v3"
)

// The configuration structure for EigenDB
type Config struct {
	Persistence struct {
		TimeInterval time.Duration `yaml:"timeInterval"`
	} `yaml:"persistence"`
	API struct {
		Port    int    `yaml:"port"`
		Address string `yaml:"address"`
	} `yaml:"api"`
	IndexConfig struct {
		Dimensions       int         `yaml:"dimensions"`
		SimilarityMetric t.SimMetric `yaml:"similarityMetric"`
	} `yaml:"indexConfig"`
}

var config *Config // the config that lives in memory

// Instantiates the in-memory config
func instantiateConfig() {
	config = new(Config)
}

// Returns a pointer to the in-memory config
func GetConfig() *Config {
	return config
}

// Writes the in-memory config to disk as a YAML file at "configPath"
//
// Returns an error if one occured.
func (c *Config) WriteToDisk(configPath string) error {
	cfgYaml, err := yaml.Marshal(config)
	if err != nil {
		return err
	}

	if err := os.WriteFile(configPath, cfgYaml, constants.CONFIG_CHMOD); err != nil {
		return err
	}
	return nil
}

// Populates the in-memory config with the values stored on disk in the YAML file at "configPath".
//
// Returns an error if one occured.
func (c *Config) populateConfig(configPath string) error {
	f, err := os.Open(configPath)
	if err != nil {
		return err
	}
	defer f.Close()

	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(config)
	if err != nil {
		return err
	}
	return nil
}

func (c *Config) populateE2EConfig() error {
	_ = c.SetPersistenceTimeInterval(3 * time.Second)
	_ = c.SetAPIPort(8080)
	_ = c.SetAPIAddress("0.0.0.0") // wouldn't 127.0.0.1 be better ?
	_ = c.SetDimensions(2)
	_ = c.SetSimilarityMetric(types.MetricL2)
	return nil
}

// Config getters and setters:
// NOTE: the setters update the specified value in-memory ONLY.

func (c *Config) GetPersistenceTimeInterval() time.Duration {
	return c.Persistence.TimeInterval
}

func (c *Config) GetAPIPort() int {
	return c.API.Port
}

func (c *Config) GetAPIAddress() string {
	return c.API.Address
}

func (c *Config) GetDimensions() int {
	return c.IndexConfig.Dimensions
}

func (c *Config) GetSimilarityMetric() t.SimMetric {
	return c.IndexConfig.SimilarityMetric
}

func (c *Config) SetPersistenceTimeInterval(timeInterval time.Duration) error {
	if timeInterval < time.Second*1 {
		return errors.New("persistence time interval must be >= 1s")
	}
	c.Persistence.TimeInterval = timeInterval
	return nil
}

func (c *Config) SetAPIPort(port int) error {
	if port <= 0 || port > 65535 {
		return errors.New("API port must be between 1 and 65535")
	}
	c.API.Port = port
	return nil
}

func (c *Config) SetAPIAddress(address string) error {
	if address == "" {
		return errors.New("API address cannot be empty")
	}
	c.API.Address = address
	return nil
}

func (c *Config) SetDimensions(dimensions int) error {
	if dimensions < 2 {
		return errors.New("dimensions must be >= 2")
	}
	c.IndexConfig.Dimensions = dimensions
	return nil
}

func (c *Config) SetSimilarityMetric(similarityMetric t.SimMetric) error {
	if err := similarityMetric.Validate(); err != nil {
		return err
	}
	c.IndexConfig.SimilarityMetric = similarityMetric
	return nil
}

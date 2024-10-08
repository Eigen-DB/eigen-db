package cfg

import (
	"os"
	"time"

	"eigen_db/constants"
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
	HNSWParams struct {
		Dimensions       int                `yaml:"dimensions"`
		SimilarityMetric t.SimilarityMetric `yaml:"similarityMetric"`
		SpaceSize        uint32             `yaml:"vectorSpaceSize"`
		M                int                `yaml:"M"`
		EfConstruction   int                `yaml:"efConstruction"`
	} `yaml:"hnswParams"`
}

var config *Config // the config that lives in memory

// Instantiates the in-memory config
func instantiateConfig() {
	config = new(Config)
}

// Returns a pointer to the in-memort config
func GetConfig() *Config {
	return config
}

// Writes the in-memory config to disk as a YAML file at "configPath"
//
// Returns an error if one occured.
func (c *Config) writeToDisk(configPath string) error {
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

// Config getters and setters:
// NOTE: the setters update the specified value in-memory AND on disk.

func (c *Config) GetPersistenceTimeInterval() time.Duration {
	return c.Persistence.TimeInterval
}

func (c *Config) GetAPIPort() int {
	return c.API.Port
}

func (c *Config) GetAPIAddress() string {
	return c.API.Address
}

func (c *Config) GetHNSWParamsDimensions() int {
	return c.HNSWParams.Dimensions
}

func (c *Config) GetHNSWParamsSimilarityMetric() t.SimilarityMetric {
	return c.HNSWParams.SimilarityMetric
}

func (c *Config) GetHNSWParamsSpaceSize() uint32 {
	return c.HNSWParams.SpaceSize
}

func (c *Config) GetHNSWParamsM() int {
	return c.HNSWParams.M
}

func (c *Config) GetHNSWParamsEfConstruction() int {
	return c.HNSWParams.EfConstruction
}

func (c *Config) SetPersistenceTimeInterval(timeInterval time.Duration) error {
	c.Persistence.TimeInterval = timeInterval
	return c.writeToDisk(constants.CONFIG_PATH)
}

func (c *Config) SetAPIPort(port int) error {
	c.API.Port = port
	return c.writeToDisk(constants.CONFIG_PATH)
}

func (c *Config) SetAPIAddress(address string) error {
	c.API.Address = address
	return c.writeToDisk(constants.CONFIG_PATH)
}

func (c *Config) SetHNSWParamsDimensions(dimensions int) error {
	c.HNSWParams.Dimensions = dimensions
	return c.writeToDisk(constants.CONFIG_PATH)
}

func (c *Config) SetHNSWParamsSimilarityMetric(similarityMetric t.SimilarityMetric) error {
	c.HNSWParams.SimilarityMetric = similarityMetric
	return c.writeToDisk(constants.CONFIG_PATH)
}

func (c *Config) SetHNSWParamsSpaceSize(spaceSize uint32) error {
	c.HNSWParams.SpaceSize = spaceSize
	return c.writeToDisk(constants.CONFIG_PATH)
}

func (c *Config) SetHNSWParamsM(M int) error {
	c.HNSWParams.M = M
	return c.writeToDisk(constants.CONFIG_PATH)
}

func (c *Config) SetHNSWParamsEfConstruction(efConstruction int) error {
	c.HNSWParams.EfConstruction = efConstruction
	return c.writeToDisk(constants.CONFIG_PATH)
}

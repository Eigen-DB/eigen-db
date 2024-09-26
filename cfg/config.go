package cfg

import (
	"os"
	"time"

	"eigen_db/constants"
	t "eigen_db/types"

	"gopkg.in/yaml.v3"
)

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

type ConfigFactory struct{}

var config *Config // the config that lives in memory

func InstantiateConfig() {
	config = new(Config)
}

func (f *ConfigFactory) GetConfig() *Config {
	return config
}

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

func (c *Config) LoadConfig(configPath string) error {
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

// Config GETTERS & SETTERS

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

func (c *Config) SetPersistenceTimeInterval(timeInterval time.Duration) {
	c.Persistence.TimeInterval = timeInterval
	c.WriteToDisk(constants.CONFIG_PATH)
}

func (c *Config) SetAPIPort(port int) {
	c.API.Port = port
	c.WriteToDisk(constants.CONFIG_PATH)
}

func (c *Config) SetAPIAddress(address string) {
	c.API.Address = address
	c.WriteToDisk(constants.CONFIG_PATH)
}

func (c *Config) SetHNSWParamsDimensions(dimensions int) {
	c.HNSWParams.Dimensions = dimensions
	c.WriteToDisk(constants.CONFIG_PATH)
}

func (c *Config) SetHNSWParamsSimilarityMetric(similarityMetric t.SimilarityMetric) {
	c.HNSWParams.SimilarityMetric = similarityMetric
	c.WriteToDisk(constants.CONFIG_PATH)
}

func (c *Config) SetHNSWParamsSpaceSize(spaceSize uint32) {
	c.HNSWParams.SpaceSize = spaceSize
	c.WriteToDisk(constants.CONFIG_PATH)
}

func (c *Config) SetHNSWParamsM(M int) {
	c.HNSWParams.M = M
	c.WriteToDisk(constants.CONFIG_PATH)
}

func (c *Config) SetHNSWParamsEfConstruction(efConstruction int) {
	c.HNSWParams.EfConstruction = efConstruction
	c.WriteToDisk(constants.CONFIG_PATH)
}

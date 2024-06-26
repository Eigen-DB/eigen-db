package cfg

import (
	"os"
	"time"

	t "eigen_db/types"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Persistence struct {
		TimeInterval time.Duration `yaml:"timeInterval"`
	} `yaml:"persistence"`
	API struct {
		Port    uint32 `yaml:"port"`
		Address string `yaml:"address"`
	} `yaml:"api"`
	HNSWParams struct {
		Dimensions       uint32             `yaml:"dimensions"`
		SimilarityMetric t.SimilarityMetric `yaml:"similarityMetric"`
		SpaceSize        uint32             `yaml:"vectorSpaceSize"`
		M                uint32             `yaml:"M"`
		EfConstruction   uint32             `yaml:"efConstruction"`
	} `yaml:"hnswParams"`
}

var config *Config

func NewConfig() {
	config = new(Config)
}

func GetConfig() *Config {
	return config
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

func (c *Config) GetAPIPort() uint32 {
	return c.API.Port
}

func (c *Config) GetAPIAddress() string {
	return c.API.Address
}

func (c *Config) GetHNSWParamsDimensions() uint32 {
	return c.HNSWParams.Dimensions
}

func (c *Config) GetHNSWParamsSimilarityMetric() t.SimilarityMetric {
	return c.HNSWParams.SimilarityMetric
}

func (c *Config) GetHNSWParamsSpaceSize() uint32 {
	return c.HNSWParams.SpaceSize
}

func (c *Config) GetHNSWParamsM() uint32 {
	return c.HNSWParams.M
}

func (c *Config) GetHNSWParamsEfConstruction() uint32 {
	return c.HNSWParams.EfConstruction
}

func (c *Config) SetPersistenceTimeInterval(timeInterval time.Duration) {
	c.Persistence.TimeInterval = timeInterval
}

func (c *Config) SetAPIPort(port uint32) {
	c.API.Port = port
}

func (c *Config) SetAPIAddress(address string) {
	c.API.Address = address
}

func (c *Config) SetHNSWParamsDimensions(dimensions uint32) {
	c.HNSWParams.Dimensions = dimensions
}

func (c *Config) SetHNSWParamsSimilarityMetric(similarityMetric t.SimilarityMetric) {
	c.HNSWParams.SimilarityMetric = similarityMetric
}

func (c *Config) SetHNSWParamsSpaceSize(spaceSize uint32) {
	c.HNSWParams.SpaceSize = spaceSize
}

func (c *Config) SetHNSWParamsM(M uint32) {
	c.HNSWParams.M = M
}

func (c *Config) SetHNSWParamsEfConstruction(efConstruction uint32) {
	c.HNSWParams.EfConstruction = efConstruction
}

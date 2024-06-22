package cfg

import (
	"eigen_db/constants"
	"os"
	"time"

	t "eigen_db/types"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Persistence struct {
		TimeIntervalSecs time.Duration `yaml:"timeIntervalSecs"`
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

var config *Config = new(Config)

func UpdateConfig() error {
	f, err := os.Open(constants.CONFIG_PATH)
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

func GetConfig() *Config {
	return config
}

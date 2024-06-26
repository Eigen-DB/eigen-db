package cfg

import (
	"eigen_db/constants"
	"fmt"
	"os"
	"testing"

	"gopkg.in/yaml.v3"
)

const CUSTOM_CONFIG_PATH string = "../" + constants.TESTING_TMP_FILES_PATH + "/custom_config.yml"

func cleanup() error {
	return os.Remove(CUSTOM_CONFIG_PATH)
}

func areAllValuesLoaded(c *Config) error {
	if c.GetPersistenceTimeInterval() == 0 {
		return fmt.Errorf("Persistence Time Interval is not set (yaml: persistence.timeInterval)")
	}
	if c.GetAPIPort() == 0 {
		return fmt.Errorf("API Port is not set (yaml: api.port)")
	}
	if c.GetAPIAddress() == "" {
		return fmt.Errorf("API Address is not set (yaml: api.address)")
	}
	if c.GetHNSWParamsDimensions() == 0 {
		return fmt.Errorf("HNSWParams Dimensions is not set (yaml: hnswParams.dimensions)")
	}
	if c.GetHNSWParamsSimilarityMetric() == "" {
		return fmt.Errorf("HNSWParams Similarity Metric is not set (yaml: hnswParams.similarityMetric)")
	}
	if c.GetHNSWParamsSpaceSize() == 0 {
		return fmt.Errorf("HNSWParams Space Size is not set (yaml: hnswParams.vectorSpaceSize)")
	}
	if c.GetHNSWParamsM() == 0 {
		return fmt.Errorf("HNSWParams M is not set (yaml: hnswParams.M)")
	}
	if c.GetHNSWParamsEfConstruction() == 0 {
		return fmt.Errorf("HNSWParams Ef Construction is not set (yaml: hnswParams.efConstruction)")
	}
	return nil
}

func areConfigsIdentical(c1 *Config, c2 *Config) error {
	if c1.GetPersistenceTimeInterval() != c2.GetPersistenceTimeInterval() {
		return fmt.Errorf("PersistenceTimeInterval values do not match. configInMem: %v, customConfigStruct: %v", c2.GetPersistenceTimeInterval(), c1.GetPersistenceTimeInterval())
	}

	if c1.GetAPIPort() != c2.GetAPIPort() {
		return fmt.Errorf("APIPort values do not match. configInMem: %v, customConfigStruct: %v", c2.GetAPIPort(), c1.GetAPIPort())
	}

	if c1.GetAPIAddress() != c2.GetAPIAddress() {
		return fmt.Errorf("APIAddress values do not match. configInMem: %v, customConfigStruct: %v", c2.GetAPIAddress(), c1.GetAPIAddress())
	}

	if c1.GetHNSWParamsDimensions() != c2.GetHNSWParamsDimensions() {
		return fmt.Errorf("HNSWParamsDimensions values do not match. configInMem: %v, customConfigStruct: %v", c2.GetHNSWParamsDimensions(), c1.GetHNSWParamsDimensions())
	}

	if c1.GetHNSWParamsSimilarityMetric() != c2.GetHNSWParamsSimilarityMetric() {
		return fmt.Errorf("HNSWParamsSimilarityMetric values do not match. configInMem: %v, customConfigStruct: %v", c2.GetHNSWParamsSimilarityMetric(), c1.GetHNSWParamsSimilarityMetric())
	}

	if c1.GetHNSWParamsSpaceSize() != c2.GetHNSWParamsSpaceSize() {
		return fmt.Errorf("HNSWParamsSpaceSize values do not match. configInMem: %v, customConfigStruct: %v", c2.GetHNSWParamsSpaceSize(), c1.GetHNSWParamsSpaceSize())
	}

	if c1.GetHNSWParamsM() != c2.GetHNSWParamsM() {
		return fmt.Errorf("HNSWParamsM values do not match. configInMem: %v, customConfigStruct: %v", c2.GetHNSWParamsM(), c1.GetHNSWParamsM())
	}

	if c1.GetHNSWParamsEfConstruction() != c2.GetHNSWParamsEfConstruction() {
		return fmt.Errorf("HNSWParamsEfConstruction values do not match. configInMem: %v, customConfigStruct: %v", c2.GetHNSWParamsEfConstruction(), c1.GetHNSWParamsEfConstruction())
	}

	return nil
}

func TestLoadConfig(t *testing.T) {
	NewConfig() // load a fresh empty config into memory
	customConfig := []byte(`
persistence:
  timeInterval: 5s
api:
  port: 8000
  address: 0.0.0.0
hnswParams:
  dimensions: 60000
  similarityMetric: euclidean
  vectorSpaceSize: 50000
  M: 15
  efConstruction: 500
`)
	if err := os.WriteFile(CUSTOM_CONFIG_PATH, customConfig, 0777); err != nil {
		t.Errorf("Error creating custom config file: %s", err.Error())
	}

	if err := GetConfig().LoadConfig(CUSTOM_CONFIG_PATH); err != nil { // load custom config into memory
		t.Errorf("Error when loading config into memory: %s", err.Error())
	}

	configInMem := GetConfig()
	if err := areAllValuesLoaded(configInMem); err != nil {
		t.Errorf(err.Error())
	}

	// check that both Config structs are identical in values
	customConfigStruct := &Config{}
	if err := yaml.Unmarshal(customConfig, customConfigStruct); err != nil {
		t.Errorf("Error parsing custom config: %s", err.Error())
	}

	if err := areConfigsIdentical(customConfigStruct, configInMem); err != nil {
		t.Errorf(err.Error())
	}

	cleanup()
}

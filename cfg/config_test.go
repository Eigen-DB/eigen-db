package cfg

import (
	"eigen_db/constants"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

const CUSTOM_CONFIG_PATH string = "../" + constants.TESTING_TMP_FILES_PATH + "/custom_config.yml"

func cleanup() error {
	return os.Remove(CUSTOM_CONFIG_PATH)
}

func areAllValuesLoaded(t *testing.T, c IConfig) {
	assert.NotEqual(t, c.GetPersistenceTimeInterval(), 0, "Persistence Time Interval is not set (yaml: persistence.timeInterval)")
	assert.NotEqual(t, c.GetAPIPort(), 0, "API Port is not set (yaml: api.port)")
	assert.NotEqual(t, c.GetAPIAddress(), "", "API Address is not set (yaml: api.address)")
	assert.NotEqual(t, c.GetHNSWParamsDimensions(), 0, "HNSWParams Dimensions is not set (yaml: hnswParams.dimensions)")
	assert.NotEqual(t, c.GetHNSWParamsSimilarityMetric(), "", "HNSWParams Similarity Metric is not set (yaml: hnswParams.similarityMetric)")
	assert.NotEqual(t, c.GetHNSWParamsSpaceSize(), 0, "HNSWParams Space Size is not set (yaml: hnswParams.vectorSpaceSize)")
	assert.NotEqual(t, c.GetHNSWParamsM(), 0, "HNSWParams M is not set (yaml: hnswParams.M)")
	assert.NotEqual(t, c.GetHNSWParamsEfConstruction(), 0, "HNSWParams Ef Construction is not set (yaml: hnswParams.efConstruction)")
}

func areConfigsIdentical(t *testing.T, c1 IConfig, c2 IConfig) {
	assert.Equal(t, c1.GetPersistenceTimeInterval(), c2.GetPersistenceTimeInterval(), "PersistenceTimeInterval values do not match. configInMem: %v, customConfigStruct: %v", c2.GetPersistenceTimeInterval(), c1.GetPersistenceTimeInterval())
	assert.Equal(t, c1.GetAPIPort(), c2.GetAPIPort(), "APIPort values do not match. configInMem: %v, customConfigStruct: %v", c2.GetAPIPort(), c1.GetAPIPort())
	assert.Equal(t, c1.GetAPIAddress(), c2.GetAPIAddress(), "APIAddress values do not match. configInMem: %v, customConfigStruct: %v", c2.GetAPIAddress(), c1.GetAPIAddress())
	assert.Equal(t, c1.GetHNSWParamsDimensions(), c2.GetHNSWParamsDimensions(), "HNSWParamsDimensions values do not match. configInMem: %v, customConfigStruct: %v", c2.GetHNSWParamsDimensions(), c1.GetHNSWParamsDimensions())
	assert.Equal(t, c1.GetHNSWParamsSimilarityMetric(), c2.GetHNSWParamsSimilarityMetric(), "HNSWParamsSimilarityMetric values do not match. configInMem: %v, customConfigStruct: %v", c2.GetHNSWParamsSimilarityMetric(), c1.GetHNSWParamsSimilarityMetric())
	assert.Equal(t, c1.GetHNSWParamsSpaceSize(), c2.GetHNSWParamsSpaceSize(), "HNSWParamsSpaceSize values do not match. configInMem: %v, customConfigStruct: %v", c2.GetHNSWParamsSpaceSize(), c1.GetHNSWParamsSpaceSize())
	assert.Equal(t, c1.GetHNSWParamsM(), c2.GetHNSWParamsM(), "HNSWParamsM values do not match. configInMem: %v, customConfigStruct: %v", c2.GetHNSWParamsM(), c1.GetHNSWParamsM())
	assert.Equal(t, c1.GetHNSWParamsEfConstruction(), c2.GetHNSWParamsEfConstruction(), "HNSWParamsEfConstruction values do not match. configInMem: %v, customConfigStruct: %v", c2.GetHNSWParamsEfConstruction(), c1.GetHNSWParamsEfConstruction())
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

	if err := (&ConfigFactory{}).GetConfig().LoadConfig(CUSTOM_CONFIG_PATH); err != nil { // load custom config into memory
		t.Errorf("Error when loading config into memory: %s", err.Error())
	}

	configInMem := (&ConfigFactory{}).GetConfig()
	areAllValuesLoaded(t, configInMem)

	// check that both Config structs are identical in values
	customConfigStruct := &Config{}
	if err := yaml.Unmarshal(customConfig, customConfigStruct); err != nil {
		t.Errorf("Error parsing custom config: %s", err.Error())
	}

	areConfigsIdentical(t, customConfigStruct, configInMem)

	cleanup()
}

// Write test for WriteToDisk method

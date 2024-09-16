package test_utils

import (
	"eigen_db/cfg"
	"eigen_db/types"
	"eigen_db/vector_io"
	"fmt"
	"time"
)

var NewVectorInvocations int = 0
var InsertInvocations int = 0
var SimilaritySearchInvocations int = 0

type MockVectorFactory struct {
	Dimensions int
}

type MockVector struct{}

type MockVectorSearcher struct{}

type MockConfigFactory struct{}

type MockConfig struct {
	WriteToDiskInvokes                   int
	LoadConfigInvokes                    int
	TimeIntervalSetInvokes               int
	TimeIntervalGetInvokes               int
	APIPortSetInvokes                    int
	APIPortGetInvokes                    int
	APIAddressSetInvokes                 int
	APIAddressGetInvokes                 int
	HNSWParamsDimensionsSetInvokes       int
	HNSWParamsDimensionsGetInvokes       int
	HNSWParamsSimilarityMetricSetInvokes int
	HNSWParamsSimilarityMetricGetInvokes int
	HNSWParamsSpaceSizeSetInvokes        int
	HNSWParamsSpaceSizeGetInvokes        int
	HNSWParamsMSetInvokes                int
	HNSWParamsMGetInvokes                int
	HNSWParamsEfConstructionSetInvokes   int
	HNSWParamsEfConstructionGetInvokes   int
}

func (factory *MockVectorFactory) NewVector(components []types.VectorComponent) (vector_io.IVector, error) {
	NewVectorInvocations++
	if len(components) == factory.Dimensions {
		return &MockVector{}, nil
	}
	return nil, fmt.Errorf("provided a %d-dimensional vector while the vector space is %d-dimensional", len(components), factory.Dimensions)
}

func (vector *MockVector) Insert() {
	InsertInvocations++
}

func (searcher *MockVectorSearcher) SimilaritySearch(queryVectorId types.VectorId, k int) ([]types.VectorId, error) {
	SimilaritySearchInvocations++
	return []types.VectorId{1, 2, 3}, nil
}

func (f *MockConfigFactory) GetConfig() cfg.IConfig {
	return &MockConfig{}
}

func (c *MockConfig) WriteToDisk(path string) error {
	c.WriteToDiskInvokes++
	return nil
}

func (c *MockConfig) LoadConfig(path string) error {
	c.LoadConfigInvokes++
	return nil
}

func (c *MockConfig) GetPersistenceTimeInterval() time.Duration {
	c.TimeIntervalGetInvokes++
	return 0
}

func (c *MockConfig) GetAPIPort() int {
	c.APIPortGetInvokes++
	return 0
}

func (c *MockConfig) GetAPIAddress() string {
	c.APIAddressGetInvokes++
	return ""
}

func (c *MockConfig) GetHNSWParamsDimensions() int {
	c.HNSWParamsDimensionsGetInvokes++
	return 0
}

func (c *MockConfig) GetHNSWParamsSimilarityMetric() types.SimilarityMetric {
	c.HNSWParamsSimilarityMetricGetInvokes++
	return types.SimilarityMetric("")
}

func (c *MockConfig) GetHNSWParamsSpaceSize() uint32 {
	c.HNSWParamsSpaceSizeGetInvokes++
	return 0
}

func (c *MockConfig) GetHNSWParamsM() int {
	c.HNSWParamsMGetInvokes++
	return 0
}

func (c *MockConfig) GetHNSWParamsEfConstruction() int {
	c.HNSWParamsEfConstructionGetInvokes++
	return 0
}

func (c *MockConfig) SetPersistenceTimeInterval(timeInterval time.Duration) error {
	c.TimeIntervalSetInvokes++
	return nil
}

func (c *MockConfig) SetAPIPort(port int) error {
	c.APIPortSetInvokes++
	return nil
}

func (c *MockConfig) SetAPIAddress(address string) error {
	c.APIAddressSetInvokes++
	return nil
}

func (c *MockConfig) SetHNSWParamsDimensions(dimensions int) error {
	c.HNSWParamsDimensionsSetInvokes++
	return nil
}

func (c *MockConfig) SetHNSWParamsSimilarityMetric(similarityMetric types.SimilarityMetric) error {
	c.HNSWParamsSimilarityMetricSetInvokes++
	return nil
}

func (c *MockConfig) SetHNSWParamsSpaceSize(spaceSize uint32) error {
	c.HNSWParamsSpaceSizeSetInvokes++
	return nil
}

func (c *MockConfig) SetHNSWParamsM(M int) error {
	c.HNSWParamsMSetInvokes++
	return nil
}

func (c *MockConfig) SetHNSWParamsEfConstruction(efConstruction int) error {
	c.HNSWParamsEfConstructionSetInvokes++
	return nil
}

func Cleanup() {
	NewVectorInvocations = 0
	InsertInvocations = 0
	SimilaritySearchInvocations = 0
}

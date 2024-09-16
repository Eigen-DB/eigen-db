package cfg

import (
	t "eigen_db/types"
	"time"
)

type IConfig interface {
	WriteToDisk(string) error
	LoadConfig(string) error

	GetPersistenceTimeInterval() time.Duration
	GetAPIPort() int
	GetAPIAddress() string
	GetHNSWParamsDimensions() int
	GetHNSWParamsSimilarityMetric() t.SimilarityMetric
	GetHNSWParamsSpaceSize() uint32
	GetHNSWParamsM() int
	GetHNSWParamsEfConstruction() int

	SetPersistenceTimeInterval(time.Duration) error
	SetAPIPort(int) error
	SetAPIAddress(string) error
	SetHNSWParamsDimensions(int) error
	SetHNSWParamsSimilarityMetric(t.SimilarityMetric) error
	SetHNSWParamsSpaceSize(uint32) error
	SetHNSWParamsM(int) error
	SetHNSWParamsEfConstruction(int) error
}

type IConfigFactory interface {
	GetConfig() IConfig
}

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

	SetPersistenceTimeInterval(time.Duration)
	SetAPIPort(int)
	SetAPIAddress(string)
	SetHNSWParamsDimensions(int)
	SetHNSWParamsSimilarityMetric(t.SimilarityMetric)
	SetHNSWParamsSpaceSize(uint32)
	SetHNSWParamsM(int)
	SetHNSWParamsEfConstruction(int)
}

type IConfigFactory interface {
	GetConfig() IConfig
}

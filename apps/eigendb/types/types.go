package types

import (
	"errors"

	"github.com/Eigen-DB/eigen-db/libs/hnswgo/v2"
)

type VecId = uint64
type Embedding = []float32
type Index = *hnswgo.Index

// Config types
type SimMetric string

const (
	COSINE        SimMetric = "cosine"
	EUCLIDEAN     SimMetric = "l2"
	INNER_PRODUCT SimMetric = "ip"
)

func (metric SimMetric) Validate() error {
	switch metric {
	case COSINE:
		return nil
	case EUCLIDEAN:
		return nil
	case INNER_PRODUCT:
		return nil
	default:
		return errors.New("invalid similarity metric")
	}
}

func (metric SimMetric) ToString() string {
	return string(metric)
}

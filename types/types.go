package types

import (
	"errors"

	"github.com/Eigen-DB/hnswgo/v2"
)

type VecId = uint64
type Embedding = []float32
type Index = *hnswgo.Index

// Config types
type SimilarityMetric string

const (
	COSINE        SimilarityMetric = "cosine"
	EUCLIDEAN     SimilarityMetric = "l2"
	INNER_PRODUCT SimilarityMetric = "ip"
)

func (metric SimilarityMetric) Validate() error {
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

func (metric SimilarityMetric) ToString() string {
	return string(metric)
}

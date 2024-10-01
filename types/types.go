package types

import (
	"errors"

	"github.com/Eigen-DB/hnswgo/v2"
)

type VectorId = uint64
type VectorComponent = float32
type Index = *hnswgo.Index

// Config types
type SimilarityMetric = string

const (
	COSINE        SimilarityMetric = "cosine"
	EUCLIDEAN     SimilarityMetric = "l2"
	INNER_PRODUCT SimilarityMetric = "ip"
)

func ParseSimilarityMetric(value string) (SimilarityMetric, error) {
	switch value {
	case COSINE:
		return COSINE, nil
	case EUCLIDEAN:
		return EUCLIDEAN, nil
	case INNER_PRODUCT:
		return INNER_PRODUCT, nil
	default:
		return "", errors.New("invalid similarity metric")
	}
}

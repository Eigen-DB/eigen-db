package types

import (
	"errors"

	"github.com/Eigen-DB/hnswgo"
)

type VectorId = int
type VectorComponent = float32
type VectorSpace = *hnswgo.HNSW

// Config types
type SimilarityMetric = string

const (
	COSINE        SimilarityMetric = "cosine"
	EUCLIDEAN     SimilarityMetric = "euclidean"
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

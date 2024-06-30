package types

import (
	"github.com/evan176/hnswgo"
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

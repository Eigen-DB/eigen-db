package types

import (
	"github.com/Eigen-DB/eigen-db/libs/faissgo"
)

type VecId = int64
type EmbeddingData = []float32
type Index = faissgo.Index

// Config types
type SimMetric = faissgo.MetricType

// func (metric SimMetric) Validate() error {
// 	switch metric {
// 	case COSINE:
// 		return nil
// 	case EUCLIDEAN:
// 		return nil
// 	case INNER_PRODUCT:
// 		return nil
// 	default:
// 		return errors.New("invalid similarity metric")
// 	}
// }

// func (metric SimMetric) ToString() string {
// 	return string(metric)
// }

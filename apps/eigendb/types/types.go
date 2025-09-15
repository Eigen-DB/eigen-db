package types

import (
	"fmt"

	"github.com/Eigen-DB/eigen-db/libs/faissgo"
)

type EmbId = int64
type EmbeddingData = []float32
type Metadata = map[string]string
type FaissIndex = faissgo.Index

// Config types
type SimMetric string

const (
	// Similarity metrics
	MetricCosine        SimMetric = "cosine"
	MetricInnerProduct  SimMetric = "ip"
	MetricL2            SimMetric = "l2"
	MetricL1            SimMetric = "l1"
	MetricLinf          SimMetric = "linf"
	MetricLp            SimMetric = "lp"
	MetricCanberra      SimMetric = "canberra"
	MetricBrayCurtis    SimMetric = "bc"
	MetricJensenShannon SimMetric = "js"
)

func (m SimMetric) String() string {
	return string(m)
}

func (m SimMetric) Validate() error {
	switch m {
	case MetricCosine, MetricInnerProduct, MetricL2, MetricL1, MetricLinf, MetricLp,
		MetricCanberra, MetricBrayCurtis, MetricJensenShannon:
		return nil
	default:
		return fmt.Errorf("invalid similarity metric: %s", m)
	}
}

func (m SimMetric) ToFaissMetricType() (faissgo.MetricType, error) {
	switch m {
	case MetricCosine, MetricInnerProduct:
		return faissgo.MetricInnerProduct, nil
	case MetricL2:
		return faissgo.MetricL2, nil
	case MetricL1:
		return faissgo.MetricL1, nil
	case MetricLinf:
		return faissgo.MetricLinf, nil
	case MetricLp:
		return faissgo.MetricLp, nil
	case MetricCanberra:
		return faissgo.MetricCanberra, nil
	case MetricBrayCurtis:
		return faissgo.MetricBrayCurtis, nil
	case MetricJensenShannon:
		return faissgo.MetricJensenShannon, nil
	default:
		return -1, fmt.Errorf("invalid similarity metric: %s", m)
	}
}

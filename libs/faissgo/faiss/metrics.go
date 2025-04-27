package faiss

/*
#cgo LDFLAGS: -lfaiss_c -lstdc++

#include <faiss/c_api/Index_c.h>
*/
import "C"

type MetricType int

const (
	MetricInnerProduct  MetricType = C.METRIC_INNER_PRODUCT
	MetricL2            MetricType = C.METRIC_L2
	MetricL1            MetricType = C.METRIC_L1
	MetricLinf          MetricType = C.METRIC_Linf
	MetricLp            MetricType = C.METRIC_Lp
	MetricCanberra      MetricType = C.METRIC_Canberra
	MetricBrayCurtis    MetricType = C.METRIC_BrayCurtis
	MetricJensenShannon MetricType = C.METRIC_JensenShannon
	//MetricJaccard         MetricType = C.METRIC_Jaccard
	//MetricNaNEuclidean    MetricType = C.METRIC_NaNEuclidean
	//MetricAbsInnerProduct MetricType = C.METRIC_ABS_INNER_PRODUCT
)

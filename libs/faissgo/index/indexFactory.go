package index

/*
#cgo CXXFLAGS: -I${SRCDIR}/../lib/faiss/c_api
#cgo LDFLAGS: -lfaiss_c -lstdc++

#include <stdlib.h>
#include <faiss/c_api/Index_c.h>
#include <faiss/c_api/index_factory_c.h>
*/
import "C"
import (
	"unsafe"

	"github.com/Eigen-DB/eigen-db/libs/faissgo/v3/faiss"
)

// desc: https://github.com/facebookresearch/faiss/wiki/The-index-factory
func IndexFactory(dim int, desc string, metric faiss.MetricType) (index, error) {
	cDesc := C.CString(desc)
	defer C.free(unsafe.Pointer(cDesc))

	idx := &faissIndex{dim: dim}
	c := C.faiss_index_factory(&idx.faissIdx, C.int(dim), cDesc, C.FaissMetricType(metric))
	if c != 0 {
		return nil, faiss.GetLastError()
	}
	return idx, nil
}

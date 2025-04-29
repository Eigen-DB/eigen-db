package faissgo

/*
#cgo LDFLAGS: -lfaiss_c -lstdc++

#include <stdlib.h>
#include <faiss/c_api/Index_c.h>
#include <faiss/c_api/index_factory_c.h>
*/
import "C"
import (
	"unsafe"
)

// desc: https://github.com/facebookresearch/faiss/wiki/The-index-factory
func IndexFactory(dim int, desc string, metric MetricType) (Index, error) {
	cDesc := C.CString(desc)
	defer C.free(unsafe.Pointer(cDesc))

	idx := &faissIndex{dim: dim}
	c := C.faiss_index_factory(&idx.faissIdx, C.int(dim), cDesc, C.FaissMetricType(metric))
	if c != 0 {
		return nil, GetLastError()
	}
	return idx, nil
}

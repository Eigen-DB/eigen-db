package faissgo

/*
#include <faiss/c_api/impl/AuxIndexStructures_c.h>
*/
import "C"

type IDSelector struct {
	sel *C.FaissIDSelector
}

func NewIDSelectorRange(idMin int64, idMax int64) (*IDSelector, error) {
	var sel *C.FaissIDSelectorRange
	c := C.faiss_IDSelectorRange_new(&sel, C.idx_t(idMin), C.idx_t(idMax))
	if c != 0 {
		return nil, GetLastError()
	}
	return &IDSelector{(*C.FaissIDSelector)(sel)}, nil
}

func NewIDSelectorBatch(ids []int64) (*IDSelector, error) {
	var sel *C.FaissIDSelectorBatch
	if c := C.faiss_IDSelectorBatch_new(
		&sel,
		C.size_t(len(ids)),
		(*C.idx_t)(&ids[0]),
	); c != 0 {
		return nil, GetLastError()
	}
	return &IDSelector{(*C.FaissIDSelector)(sel)}, nil
}

func (s *IDSelector) Free() {
	C.faiss_IDSelector_free(s.sel)
}

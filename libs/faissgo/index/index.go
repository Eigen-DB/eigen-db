package index

/*
#cgo CXXFLAGS: -I${SRCDIR}/../lib/faiss/c_api
#cgo LDFLAGS: -lfaiss_c -lstdc++

#include <stdlib.h>
#include <faiss/c_api/Index_c.h>
#include <faiss/c_api/index_io_c.h>
*/
import "C"
import (
	"unsafe"

	"github.com/Eigen-DB/eigen-db/libs/faissgo/v3/faiss"
)

type index interface {
	Train(vecsFlat []float32) error

	Add(vecsFlat []float32) error

	AddWithIds(vecsFlat []float32, ids []int64) error

	RemoveIds(n int64, ids []int64) error

	Search(queryVectsFlat []float32, k int64) ([]int64, []float32, error)

	Reconstruct(id int64) ([]float32, error)

	WriteToDisk(path string) error

	LoadFromDisk(path string) error

	IsTrained() bool

	Free()
}

type faissIndex struct {
	faissIdx *C.FaissIndex
	dim      int
}

func (idx *faissIndex) Train(vecsFlat []float32) error {
	n := int64(len(vecsFlat) / idx.dim)
	c := C.faiss_Index_train(
		idx.faissIdx,
		C.idx_t(n),
		(*C.float)(&vecsFlat[0]),
	)
	if c != 0 {
		return faiss.GetLastError()
	}
	return nil
}

func (idx *faissIndex) Add(vecsFlat []float32) error {
	n := int64(len(vecsFlat) / idx.dim)
	c := C.faiss_Index_add(
		idx.faissIdx,
		C.idx_t(n),
		(*C.float)(&vecsFlat[0]),
	)
	if c != 0 {
		return faiss.GetLastError()
	}
	return nil
}

func (idx *faissIndex) AddWithIds(vecsFlat []float32, ids []int64) error {
	n := int64(len(vecsFlat) / idx.dim)
	c := C.faiss_Index_add_with_ids(
		idx.faissIdx,
		C.idx_t(n),
		(*C.float)(&vecsFlat[0]),
		(*C.idx_t)(&ids[0]),
	)
	if c != 0 {
		return faiss.GetLastError()
	}
	return nil
}

func (idx *faissIndex) RemoveIds(n int64, ids []int64) error {
	return nil
}

func (idx *faissIndex) Search(queryVecsFlat []float32, k int64) ([]int64, []float32, error) {
	n := int64(len(queryVecsFlat) / idx.dim)
	labels := make([]int64, n*k)
	dists := make([]float32, n*k)
	c := C.faiss_Index_search(
		idx.faissIdx,
		C.idx_t(n),
		(*C.float)(&queryVecsFlat[0]),
		C.idx_t(k),
		(*C.float)(&dists[0]),
		(*C.idx_t)(&labels[0]),
	)
	if c != 0 {
		return nil, nil, faiss.GetLastError()
	}
	return labels, dists, nil
}

func (idx *faissIndex) Reconstruct(id int64) ([]float32, error) {
	v := make([]float32, idx.dim)
	c := C.faiss_Index_reconstruct(
		idx.faissIdx,
		C.idx_t(id),
		(*C.float)(&v[0]),
	)
	if c != 0 {
		return nil, faiss.GetLastError()
	}
	return v, nil
}

func (idx *faissIndex) WriteToDisk(path string) error {
	fName := C.CString(path)
	defer C.free(unsafe.Pointer(fName))
	c := C.faiss_write_index_fname(
		idx.faissIdx,
		fName,
	)
	if c != 0 {
		return faiss.GetLastError()
	}
	return nil
}

func (idx *faissIndex) LoadFromDisk(path string) error {
	fName := C.CString(path)
	defer C.free(unsafe.Pointer(fName))
	c := C.faiss_read_index_fname(
		fName,
		C.int(0), // default read only
		&idx.faissIdx,
	)
	if c != 0 {
		return faiss.GetLastError()
	}
	return nil
}

func (idx *faissIndex) IsTrained() bool {
	return C.faiss_Index_is_trained(idx.faissIdx) != 0
}

func (idx *faissIndex) Free() {
	C.faiss_Index_free(idx.faissIdx)
}

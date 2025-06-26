package faissgo

/*
#cgo LDFLAGS: -lfaiss_c -lstdc++

#include <stdlib.h>
#include <faiss/c_api/Index_c.h>
#include <faiss/c_api/impl/AuxIndexStructures_c.h>
#include <faiss/c_api/index_io_c.h>
*/
import "C"
import (
	"unsafe"
)

type Index interface {
	Train(vecsFlat []float32) error

	Add(vecsFlat []float32) error

	AddWithIds(vecsFlat []float32, ids []int64) error

	RemoveIds(selector *IDSelector) (int, error)

	Search(queryVectsFlat []float32, k int64) ([]int64, []float32, error)

	Reconstruct(id int64) ([]float32, error) // can be used to fetch a vector for any type of index (i think)

	ReconstructN(id int64, n int64) ([]float32, error) // Reconstruct vectors id to id + n - 1

	WriteToDisk(path string) error

	LoadFromDisk(path string) error

	IsTrained() bool

	NTotal() int64

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
		return GetLastError()
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
		return GetLastError()
	}
	return nil
}

// if ID already exists, it will be overwritten, but the old vector will not be removed, but just marked as deleted.
// this means that NTotal() will still return the total number of vectors added to the index, even if some of them are marked as deleted.
func (idx *faissIndex) AddWithIds(vecsFlat []float32, ids []int64) error {
	n := int64(len(vecsFlat) / idx.dim)
	c := C.faiss_Index_add_with_ids(
		idx.faissIdx,
		C.idx_t(n),
		(*C.float)(&vecsFlat[0]),
		(*C.idx_t)(&ids[0]),
	)
	if c != 0 {
		return GetLastError()
	}
	return nil
}

func (idx *faissIndex) RemoveIds(selector *IDSelector) (int, error) {
	var nRemoved C.size_t
	c := C.faiss_Index_remove_ids(
		idx.faissIdx,
		selector.sel,
		&nRemoved,
	)
	if c != 0 {
		return 0, GetLastError()
	}
	return int(nRemoved), nil
}

// supports many query vectors at once
func (idx *faissIndex) Search(queryVecsFlat []float32, k int64) ([]int64, []float32, error) {
	n := int64(len(queryVecsFlat) / idx.dim) // number of query vectors
	var labels []int64
	var dists []float32
	if k > idx.NTotal() {
		labels = make([]int64, n*idx.NTotal())
		dists = make([]float32, n*idx.NTotal())
	} else {
		labels = make([]int64, n*k)
		dists = make([]float32, n*k)
	}
	c := C.faiss_Index_search(
		idx.faissIdx,
		C.idx_t(n),
		(*C.float)(&queryVecsFlat[0]),
		C.idx_t(k),
		(*C.float)(&dists[0]),
		(*C.idx_t)(&labels[0]),
	)
	if c != 0 {
		return nil, nil, GetLastError()
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
		return nil, GetLastError()
	}
	return v, nil
}

func (idx *faissIndex) ReconstructN(id int64, n int64) ([]float32, error) {
	// TODO: prevent function caller from making the range of reconstructed vector out of bounds of the index
	// simple solution: if id + n > idx.NTotal() -> truncate the range to fit within the bound of the index
	// problem: this solution works when the index type requires vectors to be added without an ID (increasing from 0 to n-1).
	// It is still unclear how this function will behave with an index with custom vector IDs.
	vecs := make([]float32, n*int64(idx.dim))
	c := C.faiss_Index_reconstruct_n(
		idx.faissIdx,
		C.idx_t(id),
		C.idx_t(n),
		(*C.float)(&vecs[0]),
	)
	if c != 0 {
		return nil, GetLastError()
	}
	return vecs, nil
}

func (idx *faissIndex) WriteToDisk(path string) error {
	fName := C.CString(path)
	defer C.free(unsafe.Pointer(fName))
	c := C.faiss_write_index_fname(
		idx.faissIdx,
		fName,
	)
	if c != 0 {
		return GetLastError()
	}
	return nil
}

// When loading an index from disk, the index description passed into the Index Factory
// doesn't matter as the index you created to load the index from disk will be replaced
// by the one on disk.
func (idx *faissIndex) LoadFromDisk(path string) error {
	fName := C.CString(path)
	defer C.free(unsafe.Pointer(fName))
	c := C.faiss_read_index_fname(
		fName,
		C.int(0), // default read only
		&idx.faissIdx,
	)
	if c != 0 {
		return GetLastError()
	}

	return nil
}

func (idx *faissIndex) IsTrained() bool {
	return C.faiss_Index_is_trained(idx.faissIdx) != 0
}

func (idx *faissIndex) NTotal() int64 {
	n := C.faiss_Index_ntotal(idx.faissIdx)
	// if n < 0 {
	// 	return 0
	// }
	return int64(n)
}

func (idx *faissIndex) Free() {
	C.faiss_Index_free(idx.faissIdx)
}

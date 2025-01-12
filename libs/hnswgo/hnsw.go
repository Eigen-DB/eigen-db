package hnswgo

/*
#cgo CXXFLAGS: -std=c++11
#include <stdlib.h>
#include <hnsw_wrapper.h>
*/
import "C"
import (
	"errors"
	"fmt"
	"math"
	"os"
	"unsafe"
)

type Index struct {
	index      C.HNSW
	dimensions int
	size       uint32
	normalize  bool
	spaceType  string
}

// Returns the last error message. Returns nil if there is no error message.
func peekLastError() error {
	err := C.peekLastErrorMsg()
	if err == nil {
		return nil
	}
	return errors.New(C.GoString(err))
}

// Returns and clears the last error message. Returns nil if there is no error message.
func getLastError() error {
	err := C.getLastErrorMsg()
	if err == nil {
		return nil
	}
	return errors.New(C.GoString(err))
}

/*
Normalizes a vector in place.
Normalize(v) = (1/|v|)*v

- vector: the vector to Normalize in place
*/
func Normalize(vector []float32) {
	var magnitude float32
	for i := range vector {
		magnitude += vector[i] * vector[i]
	}
	magnitude = float32(math.Sqrt(float64(magnitude)))

	for i := range vector {
		vector[i] *= 1.0 / magnitude
	}
}

/*
Returns a reference to an instance of an HNSW index.

- dim:            	dimension of the vector space

- maxElements:    	index's vector storage capacity

- m:              	`m` parameter in the HNSW algorithm

- efConstruction: 	`efConstruction` parameter in the HNSW algorithm

- randSeed:       	random seed

- spaceType:      	similarity metric to use in the index ("ip", "cosine", "l2". default: "l2")

Returns an instance of an HNSW index, or an error if there was a problem initializing the index.
*/
func New(dim int, m int, efConstruction int, randSeed int, maxElements uint32, spaceType string) (*Index, error) {
	if dim < 1 {
		return nil, errors.New("dimension must be >= 1")
	}
	if maxElements < 1 {
		return nil, errors.New("max elements must be >= 1")
	}
	if m < 2 {
		return nil, errors.New("m must be >= 2")
	}
	if efConstruction < 0 {
		return nil, errors.New("efConstruction must be >= 0")
	}

	index := new(Index)
	index.dimensions = dim
	index.spaceType = spaceType
	index.size = maxElements

	if spaceType == "ip" {
		index.index = C.initHNSW(C.int(dim), C.ulong(maxElements), C.int(m), C.int(efConstruction), C.int(randSeed), C.char('i'))
	} else if spaceType == "cosine" {
		index.normalize = true
		index.index = C.initHNSW(C.int(dim), C.ulong(maxElements), C.int(m), C.int(efConstruction), C.int(randSeed), C.char('c'))
	} else {
		index.index = C.initHNSW(C.int(dim), C.ulong(maxElements), C.int(m), C.int(efConstruction), C.int(randSeed), C.char('l'))
	}

	if index.index == nil {
		return nil, getLastError()
	}

	return index, getLastError()
}

/*
Loads a saved index and returns a reference to it.

- location:			the file path of the saved index

- dim:            	dimension of the vector space

- spaceType:      	similarity metric to use in the index ("ip", "cosine", "l2". default: "l2")

- maxElements:    	index's vector storage capacity

Returns an instance of the saved HNSW index, or an error if there was a problem.
*/
func LoadIndex(location string, dim int, spaceType string, maxElements uint32) (*Index, error) {
	if dim < 1 {
		return nil, errors.New("dimension must be >= 1")
	}
	if maxElements < 1 {
		return nil, errors.New("max elements must be >= 1")
	}

	// checking the location's validity and permissions
	if _, err := os.ReadFile(location); err != nil {
		return nil, err
	}

	index := new(Index)
	index.dimensions = dim
	index.spaceType = spaceType
	index.size = maxElements

	cLocation := C.CString(location)
	defer C.free(unsafe.Pointer(cLocation))

	if spaceType == "ip" {
		index.index = C.loadHNSW(cLocation, C.int(dim), C.char('i'), C.ulong(maxElements))
	} else if spaceType == "cosine" {
		index.normalize = true
		index.index = C.loadHNSW(cLocation, C.int(dim), C.char('c'), C.ulong(maxElements))
	} else {
		index.index = C.loadHNSW(cLocation, C.int(dim), C.char('l'), C.ulong(maxElements))
	}

	if index.index == nil {
		return nil, getLastError()
	}

	return index, getLastError()
}

/*
Saves the index to the disk.

- location:			the file path in which to save the index

Returns an error if there was a problem.
*/
func (i *Index) SaveToDisk(location string) error {
	if location == "" {
		return errors.New("location cannot be blank")
	}

	// checking the location's validity and permissions
	if _, err := os.Stat(location); os.IsNotExist(err) { // file does not exist yet
		if _, err := os.Create(location); err != nil {
			return err
		}
		if err := os.Remove(location); err != nil {
			return err
		}
	} else if os.IsPermission(err) {
		return err
	}

	cLocation := C.CString(location)
	defer C.free(unsafe.Pointer(cLocation))
	C.saveHNSW(i.index, cLocation)
	return getLastError()
}

/*
Frees the HNSW index from memory.
*/
func (i *Index) Free() {
	C.freeHNSW(i.index)
}

/*
Adds a vector to the HNSW index.
If the a vector with the same label already exists, the function returns an error

- vector:       the vector to add to the index

- label:        the vector's label

Returns an error if one occured.
*/
func (i *Index) InsertVector(vector []float32, label uint64) error {
	if len(vector) != i.dimensions {
		return fmt.Errorf("the vector you are trying to insert is %d-dimensional whereas your index is %d-dimensional", len(vector), i.dimensions)
	}

	_, err := i.GetVector(label)
	if err == nil {
		return fmt.Errorf("a vector with label %d already exists in the index", label)
	}

	if i.normalize {
		Normalize(vector)
	}
	C.insertVector(i.index, (*C.float)(unsafe.Pointer(&vector[0])), C.ulong(label))
	return getLastError()
}

/*
Replaces an existing vector in the HNSW index.

- label:        the vector's label

- newVector:    the new vector used to replace the old vector

Returns an error if one occured.
*/
func (i *Index) ReplaceVector(label uint64, newVector []float32) error {
	if len(newVector) != i.dimensions {
		return fmt.Errorf("the vector you are trying to insert is %d-dimensional whereas your index is %d-dimensional", len(newVector), i.dimensions)
	}
	if i.normalize {
		Normalize(newVector)
	}
	C.insertVector(i.index, (*C.float)(unsafe.Pointer(&newVector[0])), C.ulong(label))
	return getLastError()
}

/*
Returns a vector's components using its label

- label:	the vector's label

Returns the vector's components with specified label
*/
func (i *Index) GetVector(label uint64) ([]float32, error) {
	vec := C.getVector(i.index, C.ulong(label), C.int(i.dimensions))
	if vec == nil {
		return nil, getLastError()
	}
	defer C.free(unsafe.Pointer(vec))
	vecSlice := make([]float32, i.dimensions)
	copy(vecSlice, unsafe.Slice((*float32)(vec), i.dimensions))

	return vecSlice, getLastError()
}

/*
Marks a vector with the specified label as deleted, which omits it from KNN search.

- label: 	the vector's label

Returns an error if one occured.
*/
func (i *Index) DeleteVector(label uint64) error {
	C.deleteVector(i.index, C.ulong(label))
	return getLastError()
}

/*
Performs similarity search on the HNSW index.

- vector:       the query vector

- k:            the k value

Returns the labels and distances of each of the nearest neighbors, and an error if one occured. Note: the size of both arrays can be < k if k > num of vectors in the index
*/
func (i *Index) SearchKNN(vector []float32, k int) ([]uint64, []float32, error) {
	if len(vector) != i.dimensions {
		return nil, nil, fmt.Errorf("the query vector is %d-dimensional whereas your index is %d-dimensional", len(vector), i.dimensions)
	}
	if k < 1 || uint32(k) > i.size {
		return nil, nil, fmt.Errorf("1 <= k <= index max size")
	}

	if i.normalize {
		Normalize(vector)
	}

	Clabel := make([]C.ulong, k)
	Cdist := make([]C.float, k)

	numResult := int(C.searchKNN(i.index, (*C.float)(unsafe.Pointer(&vector[0])), C.int(k), &Clabel[0], &Cdist[0])) // perform the search

	if numResult < 0 {
		return nil, nil, fmt.Errorf("an error occured with the HNSW algorithm: %s", getLastError())
	}

	labels := make([]uint64, k)
	dists := make([]float32, k)
	for i := 0; i < numResult; i++ {
		labels[i] = uint64(Clabel[i])
		dists[i] = float32(Cdist[i])
	}

	return labels[:numResult], dists[:numResult], getLastError()
}

/*
Set's the efConstruction parameter in the HNSW index.

- efConstruction: the new efConstruction parameter

Returns an error if one occured.
*/
func (i *Index) SetEfConstruction(efConstruction int) error {
	if efConstruction < 0 {
		return errors.New("efConstruction must be >= 0")
	}
	C.setEf(i.index, C.int(efConstruction))
	return getLastError()
}

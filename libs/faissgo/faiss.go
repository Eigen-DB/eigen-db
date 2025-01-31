package faissgo

import "C"
import (
	"math"
)

type Index struct {
}

// Returns the last error message. Returns nil if there is no error message.
func peekLastError() error {
	return nil
}

// Returns and clears the last error message. Returns nil if there is no error message.
func getLastError() error {
	return nil
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
	return nil, nil
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
	return nil, nil
}

/*
Saves the index to the disk.

- location:			the file path in which to save the index

Returns an error if there was a problem.
*/
func (i *Index) SaveToDisk(location string) error {
	return nil
}

/*
Frees the HNSW index from memory.
*/
func (i *Index) Free() {

}

/*
Adds a vector to the HNSW index.
If the a vector with the same label already exists, the function returns an error

- vector:       the vector to add to the index

- label:        the vector's label

Returns an error if one occured.
*/
func (i *Index) InsertVector(vector []float32, label uint64) error {
	return nil
}

/*
Replaces an existing vector in the HNSW index.

- label:        the vector's label

- newVector:    the new vector used to replace the old vector

Returns an error if one occured.
*/
func (i *Index) ReplaceVector(label uint64, newVector []float32) error {
	return nil
}

/*
Returns a vector's components using its label

- label:	the vector's label

Returns the vector's components with specified label
*/
func (i *Index) GetVector(label uint64) ([]float32, error) {
	return nil, nil
}

/*
Marks a vector with the specified label as deleted, which omits it from KNN search.

- label: 	the vector's label

Returns an error if one occured.
*/
func (i *Index) DeleteVector(label uint64) error {
	return nil
}

/*
Performs similarity search on the HNSW index.

- vector:       the query vector

- k:            the k value

Returns the labels and distances of each of the nearest neighbors, and an error if one occured. Note: the size of both arrays can be < k if k > num of vectors in the index
*/
func (i *Index) SearchKNN(vector []float32, k int) ([]uint64, []float32, error) {
	return nil, nil, nil
}

/*
Set's the efConstruction parameter in the HNSW index.

- efConstruction: the new efConstruction parameter

Returns an error if one occured.
*/
func (i *Index) SetEfConstruction(efConstruction int) error {
	return nil
}

package tests

import (
	"fmt"
	"math"
	"os"
	"sort"
	"testing"
	"time"

	"github.com/Eigen-DB/eigen-db/libs/hnswgo/v2"

	"github.com/stretchr/testify/assert"
)

func setup() (*hnswgo.Index, error) {
	index, err := hnswgo.New(
		2,
		32,
		400,
		int(time.Now().Unix()),
		uint32(10000),
		"l2",
	)

	if err != nil {
		return nil, fmt.Errorf("An error occured when instantiating the index: %s", err.Error())
	}
	return index, nil
}

func TestNormalize(t *testing.T) {
	vector := []float32{4, 5, 6}
	hnswgo.Normalize(vector)
	var magnitude float32
	for i := range vector {
		magnitude += vector[i] * vector[i]
	}
	magnitude = float32(math.Sqrt(float64(magnitude)))
	assert.Equal(t, magnitude, float32(1.0))
}

func TestNewSuccess(t *testing.T) {
	dim := 2
	maxElements := uint32(10000)
	m := 32
	efConstruction := 400
	spaceType := "l2"
	seed := int(time.Now().Unix())

	index, err := hnswgo.New(
		dim,
		m,
		efConstruction,
		seed,
		maxElements,
		spaceType,
	)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	} else {
		defer index.Free()
	}

	if index == nil {
		t.Fatal("Expected valid index, got nil")
	}
}

func TestNewFailure(t *testing.T) {
	dim := -128 // Invalid dimension
	m := 16
	efConstruction := 200
	randSeed := int(time.Now().Unix())
	maxElements := uint32(10000)
	spaceType := "l2"

	index, err := hnswgo.New(
		dim,
		m,
		efConstruction,
		randSeed,
		maxElements,
		spaceType,
	)

	if err == nil {
		t.Fatal("Expected an error for invalid parameters, but got none")
	}

	if index != nil {
		defer index.Free()
		t.Fatal("Expected nil index on failure, but got valid index")
	}
}

func TestInsertVectorSuccess(t *testing.T) {
	index, err := setup()
	if err != nil {
		t.Fatal(err.Error())
	}
	defer index.Free()

	vector := []float32{1.2, -4.2}
	if err := index.InsertVector(vector, 1); err != nil {
		t.Fatalf("An error occured when inserting a vector: %s", err.Error())
	}
}

func TestInsertVectorFailure(t *testing.T) {
	index, err := setup()
	if err != nil {
		t.Fatal(err.Error())
	}
	defer index.Free()

	vector := []float32{1.2, -4.2, 3.3} // trying to insert 3-dimensional vector in 2-dimensional index -> error
	if err := index.InsertVector(vector, 1); err == nil {
		t.Fatal("An error SHOULD HAVE occured when inserting a 3D vector in a 2D index")
	}
}

func TestInsertVectorOverwrite(t *testing.T) {
	index, err := setup()
	if err != nil {
		t.Fatal(err.Error())
	}
	defer index.Free()

	if err := index.InsertVector([]float32{1.2, -4.2}, 1); err != nil {
		t.Fatalf("An error occured when inserting a vector: %s", err.Error())
	}

	if err := index.InsertVector([]float32{4.2, 6.2}, 1); err == nil {
		t.Fatalf("No error occured when attempting to overwrite a vector")
	} else if err.Error() != "a vector with label 1 already exists in the index" {
		t.Fatalf("Got the wrong error when trying to overwrite a vector: %s", err.Error())
	}
}

func TestReplaceVectorSuccess(t *testing.T) {
	index, err := setup()
	if err != nil {
		t.Fatal(err.Error())
	}
	defer index.Free()

	if err := index.InsertVector([]float32{1.2, -4.2}, 1); err != nil {
		t.Fatalf("An error occured when inserting a vector: %s", err.Error())
	}

	if err := index.ReplaceVector(1, []float32{4.2, 6.2}); err != nil {
		t.Fatalf("An error occured when trying to replace a vector: %s", err.Error())
	}
}

func TestReplaceVectorFirstInsert(t *testing.T) {
	index, err := setup()
	if err != nil {
		t.Fatal(err.Error())
	}
	defer index.Free()

	if err := index.ReplaceVector(1, []float32{4.2, 6.2}); err != nil {
		t.Fatalf("An error occured when trying to replace a non-existant vector: %s", err.Error())
	}
}

func TestReplaceVectorInvalidDims(t *testing.T) {
	index, err := setup()
	if err != nil {
		t.Fatal(err.Error())
	}
	defer index.Free()

	if err := index.InsertVector([]float32{1.2, -4.2}, 1); err != nil {
		t.Fatalf("An error occured when inserting a vector: %s", err.Error())
	}

	if err := index.ReplaceVector(1, []float32{4.2, 6.2, 3.3}); err == nil {
		t.Fatalf("No error occured when trying to replace a vector with invalid dimensions")
	} else if err.Error() != "the vector you are trying to insert is 3-dimensional whereas your index is 2-dimensional" {
		t.Fatalf("Got the wrong error when trying to overwrite a vector with invalid dimensions: %s", err.Error())
	}
}

func TestGetVectorSuccess(t *testing.T) {
	index, err := setup()
	if err != nil {
		t.Fatal(err.Error())
	}
	defer index.Free()

	// sample vectors
	vectors := [][]float32{
		{1.2, 3.4},
		{2.1, 4.5},
		{0.5, 1.7},
		{3.3, 2.2},
		{4.8, 5.6},
		{7.1, 8.2},
		{9.0, 0.4},
		{6.3, 3.5},
		{2.9, 7.8},
		{5.0, 1.1},
	}

	// insert sample vectors
	for i, v := range vectors {
		if err := index.InsertVector(v, uint64(i)); err != nil {
			t.Fatalf("An error occured when inserting vector %v: %s", v, err.Error())
		}
	}

	vec, err := index.GetVector(2)
	if err != nil {
		t.Fatalf("An error occured when getting a vector: %s", err.Error())
	}

	assert.Equal(t, vec, vectors[2], fmt.Sprintf("vector gotten != expected vector. Vector gotten: %v. Expected: %v.", vec, vectors[2]))
}

func TestGetVectorNotFound(t *testing.T) {
	index, err := setup()
	if err != nil {
		t.Fatal(err.Error())
	}
	defer index.Free()

	// sample vectors
	vectors := [][]float32{
		{1.2, 3.4},
		{2.1, 4.5},
		{0.5, 1.7},
		{3.3, 2.2},
		{4.8, 5.6},
		{7.1, 8.2},
		{9.0, 0.4},
		{6.3, 3.5},
		{2.9, 7.8},
		{5.0, 1.1},
	}

	// insert sample vectors
	for i, v := range vectors {
		if err := index.InsertVector(v, uint64(i)); err != nil {
			t.Fatalf("An error occured when inserting vector %v: %s", v, err.Error())
		}
	}

	_, err = index.GetVector(100)
	if err == nil {
		t.Fatal("No error occured when getting a non-existant vector")
	} else {
		if err.Error() != "Label not found" {
			t.Fatalf("Another error OTHER than \"Label not found\" occured: %s", err.Error())
		}
	}
}

func TestDeleteVectorSuccessful(t *testing.T) {
	index, err := setup()
	if err != nil {
		t.Fatal(err.Error())
	}
	defer index.Free()

	// sample vectors
	vectors := [][]float32{
		{1.2, 3.4},
		{2.1, 4.5},
		{0.5, 1.7},
		{3.3, 2.2},
		{4.8, 5.6},
		{7.1, 8.2},
		{9.0, 0.4},
		{6.3, 3.5},
		{2.9, 7.8},
		{5.0, 1.1},
	}

	// insert sample vectors
	for i, v := range vectors {
		if err := index.InsertVector(v, uint64(i)); err != nil {
			t.Fatalf("An error occured when inserting vector %v: %s", v, err.Error())
		}
	}

	if err := index.DeleteVector(2); err != nil {
		t.Fatalf("An error occured when trying to delete a vector: %s", err.Error())
	}

	if _, err = index.GetVector(2); err != nil {
		if err.Error() != "Label not found" {
			t.Fatalf("An error OTHER than \"Label not found\" occured when trying to get a vector AFTER it's been deleted")
		}
	} else {
		t.Fatalf("No error occured when trying to get a vector AFTER it's been deleted")
	}

	if err := index.InsertVector(vectors[2], uint64(2)); err != nil {
		t.Fatalf("An error occured when trying to re-insert the same vector with the same label after deletion: %s", err.Error())
	}

	v, err := index.GetVector(2)
	if err != nil {
		t.Fatalf("An error occured when trying to get a vector after it's been re-instered: %s", err.Error())
	}

	assert.Equal(t, v, vectors[2], fmt.Sprintf("The vector you re-inserted is not the same as the actual vector. Vector re-inserted: %v. Expected: %v.", v, vectors[2]))
}

func TestDeleteVectorUnsuccessful(t *testing.T) {
	index, err := setup()
	if err != nil {
		t.Fatal(err.Error())
	}
	defer index.Free()

	// sample vectors
	vectors := [][]float32{
		{1.2, 3.4},
		{2.1, 4.5},
		{0.5, 1.7},
		{3.3, 2.2},
		{4.8, 5.6},
		{7.1, 8.2},
		{9.0, 0.4},
		{6.3, 3.5},
		{2.9, 7.8},
		{5.0, 1.1},
	}

	// insert sample vectors
	for i, v := range vectors {
		if err := index.InsertVector(v, uint64(i)); err != nil {
			t.Fatalf("An error occured when inserting vector %v: %s", v, err.Error())
		}
	}

	if err := index.DeleteVector(100); err == nil {
		t.Fatal("No error occured when trying to delete a non-existant vector")
	} else {
		if err.Error() != "Label not found" {
			t.Fatalf("An error other than \"Label not found\" occured when trying to delete a non-existant vector")
		}
	}
}

func TestDeleteVectorKNNAfterReinsertion(t *testing.T) {
	index, err := setup()
	if err != nil {
		t.Fatal(err.Error())
	}
	defer index.Free()

	// sample vectors
	vectors := [][]float32{
		{1.2, 3.4},
		{2.1, 4.5},
		{0.5, 1.7},
		{3.3, 2.2},
		{4.8, 5.6},
		{7.1, 8.2},
		{9.0, 0.4},
		{6.3, 3.5},
		{2.9, 7.8},
		{5.0, 1.1},
	}

	// insert sample vectors
	for i, v := range vectors {
		if err := index.InsertVector(v, uint64(i)); err != nil {
			t.Fatalf("An error occured when inserting vector %v: %s", v, err.Error())
		}
	}

	if err := index.DeleteVector(2); err != nil {
		t.Fatalf("An error occured when trying to delete a vector: %s", err.Error())
	}

	if _, err = index.GetVector(2); err != nil {
		if err.Error() != "Label not found" {
			t.Fatalf("An error OTHER than \"Label not found\" occured when trying to get a vector AFTER it's been deleted")
		}
	} else {
		t.Fatalf("No error occured when trying to get a vector AFTER it's been deleted")
	}

	if err := index.InsertVector(vectors[2], uint64(2)); err != nil {
		t.Fatalf("An error occured when trying to re-insert the same vector with the same label after deletion: %s", err.Error())
	}

	k := 5
	nnLabels, _, err := index.SearchKNN(vectors[0], k)
	if err != nil {
		t.Fatalf("Error when performing similarity search: %s", err.Error())
	}

	sort.Slice(nnLabels, func(i, j int) bool {
		return nnLabels[i] < nnLabels[j]
	})

	assert.Equal(t, []uint64{0, 1, 2, 3, 4}, nnLabels)
}

func TestSearchKNN(t *testing.T) {
	index, err := setup()
	if err != nil {
		t.Fatal(err.Error())
	}
	defer index.Free()

	// sample vectors
	vectors := [][]float32{
		{1.2, 3.4},
		{2.1, 4.5},
		{0.5, 1.7},
		{3.3, 2.2},
		{4.8, 5.6},
		{7.1, 8.2},
		{9.0, 0.4},
		{6.3, 3.5},
		{2.9, 7.8},
		{5.0, 1.1},
	}

	// insert sample vectors
	for i, v := range vectors {
		if err := index.InsertVector(v, uint64(i)); err != nil {
			t.Fatalf("An error occured when inserting vector %v: %s", v, err.Error())
		}
	}

	k := 5
	nnLabels, nnDists, err := index.SearchKNN(vectors[0], k) // perform similarity search where the first of our sample vectors is the query vector
	if err != nil {
		t.Fatalf("Error when performing similarity search: %s", err.Error())
	}

	sort.Slice(nnLabels, func(i, j int) bool {
		return nnLabels[i] < nnLabels[j]
	})

	assert.Equal(t, []uint64{0, 1, 2, 3, 4}, nnLabels)

	t.Logf("%d-nearest neighbors:\n", k)
	for i := range nnLabels {
		t.Logf("vector %d is %f units from query vector\n", nnLabels[i], nnDists[i])
	}
}

func TestSetEfConstructionSuccess(t *testing.T) {
	index, err := setup()
	if err != nil {
		t.Fatal(err.Error())
	}
	defer index.Free()

	err = index.SetEfConstruction(401)
	if err != nil {
		t.Fatalf("An error occured when updating efConstruction: %s", err.Error())
	}
}

func TestSetEfConstructionFailure(t *testing.T) {
	index, err := setup()
	if err != nil {
		t.Fatal(err.Error())
	}
	defer index.Free()

	err = index.SetEfConstruction(-1)
	if err == nil {
		t.Fatal("An error SHOULD HAVE occured when updating efConstruction.")
	}
}

func TestSaveToDiskSuccess(t *testing.T) {
	index, err := setup()
	if err != nil {
		t.Fatal(err.Error())
	}
	defer index.Free()

	// sample vectors
	vectors := [][]float32{
		{1.2, 3.4},
		{2.1, 4.5},
		{0.5, 1.7},
		{3.3, 2.2},
		{4.8, 5.6},
		{7.1, 8.2},
		{9.0, 0.4},
		{6.3, 3.5},
		{2.9, 7.8},
		{5.0, 1.1},
	}

	// insert sample vectors
	for i, v := range vectors {
		if err := index.InsertVector(v, uint64(i)); err != nil {
			t.Fatalf("An error occured when inserting vector %v: %s", v, err.Error())
		}
	}

	if err := index.SaveToDisk("/tmp/saved_data.bin"); err != nil {
		t.Fatalf("An error occured when saving index to disk: %s", err.Error())
	}
}

func TestSaveToDiskEmptyPath(t *testing.T) {
	index, err := setup()
	if err != nil {
		t.Fatal(err.Error())
	}
	defer index.Free()

	// sample vectors
	vectors := [][]float32{
		{1.2, 3.4},
		{2.1, 4.5},
		{0.5, 1.7},
		{3.3, 2.2},
		{4.8, 5.6},
		{7.1, 8.2},
		{9.0, 0.4},
		{6.3, 3.5},
		{2.9, 7.8},
		{5.0, 1.1},
	}

	// insert sample vectors
	for i, v := range vectors {
		if err := index.InsertVector(v, uint64(i)); err != nil {
			t.Fatalf("An error occured when inserting vector %v: %s", v, err.Error())
		}
	}

	if err := index.SaveToDisk(""); err == nil {
		t.Fatal("No error occured when setting the location to \"\"")
	} else if err.Error() != "location cannot be blank" {
		t.Fatalf("Unexpected error received: %s", err.Error())
	}
}

func TestSaveToDiskInvalidPath(t *testing.T) {
	index, err := setup()
	if err != nil {
		t.Fatal(err.Error())
	}
	defer index.Free()

	// sample vectors
	vectors := [][]float32{
		{1.2, 3.4},
		{2.1, 4.5},
		{0.5, 1.7},
		{3.3, 2.2},
		{4.8, 5.6},
		{7.1, 8.2},
		{9.0, 0.4},
		{6.3, 3.5},
		{2.9, 7.8},
		{5.0, 1.1},
	}

	// insert sample vectors
	for i, v := range vectors {
		if err := index.InsertVector(v, uint64(i)); err != nil {
			t.Fatalf("An error occured when inserting vector %v: %s", v, err.Error())
		}
	}

	if err := index.SaveToDisk("/some/fake/path/index.bin"); err == nil {
		t.Fatal("No error occured when setting the location to an invalid path")
	} else if !os.IsNotExist(err) {
		t.Fatalf("Unexpected error received: %s", err.Error())
	}
}

func TestSaveToDiskInvalidPerms(t *testing.T) {
	index, err := setup()
	if err != nil {
		t.Fatal(err.Error())
	}
	defer index.Free()

	// sample vectors
	vectors := [][]float32{
		{1.2, 3.4},
		{2.1, 4.5},
		{0.5, 1.7},
		{3.3, 2.2},
		{4.8, 5.6},
		{7.1, 8.2},
		{9.0, 0.4},
		{6.3, 3.5},
		{2.9, 7.8},
		{5.0, 1.1},
	}

	// insert sample vectors
	for i, v := range vectors {
		if err := index.InsertVector(v, uint64(i)); err != nil {
			t.Fatalf("An error occured when inserting vector %v: %s", v, err.Error())
		}
	}

	if err := index.SaveToDisk("/root/index.bin"); err == nil {
		t.Fatal("No error occured when setting the location to a path with elevated permissions")
	} else if !os.IsPermission(err) {
		t.Fatalf("Unexpected error received: %s", err.Error())
	}
}

func TestLoadIndexSuccess(t *testing.T) {
	index, err := hnswgo.LoadIndex("/tmp/saved_data.bin", 2, "l2", uint32(10000)) // same as whats in setup()
	if err != nil {
		t.Fatalf("An error occured when loading an index from disk: %s", err.Error())
	}
	defer index.Free()

	expectedVectors := [][]float32{
		{1.2, 3.4},
		{2.1, 4.5},
		{0.5, 1.7},
		{3.3, 2.2},
		{4.8, 5.6},
		{7.1, 8.2},
		{9.0, 0.4},
		{6.3, 3.5},
		{2.9, 7.8},
		{5.0, 1.1},
	}

	for i := 0; i < 10; i++ {
		v, err := index.GetVector(uint64(i))
		if err != nil {
			t.Fatalf("An error occured when getting a vector from the index loaded from disk: %s", err.Error())
		}
		assert.Equal(t, expectedVectors[i], v, fmt.Sprintf("The vector fetched from loaded index != vector inserted before saving on disk. Expected: %v. Got: %v.", expectedVectors[i], v))

		t.Logf("%v\n", v)
	}

	// cleanup
	if err := os.Remove("/tmp/saved_data.bin"); err != nil {
		t.Fatal(err.Error())
	}
}

func TestLoadIndexInvalidPath(t *testing.T) {
	index, err := hnswgo.LoadIndex("/fake/path", 2, "l2", uint32(10000))
	if index != nil {
		defer index.Free()
	}
	if err == nil {
		t.Fatal("No error occured when trying to load an index with a fake location.")
	} else if !os.IsNotExist(err) {
		t.Fatalf("Unexpected error received: %s", err.Error())
	}
}

func TestLoadIndexEmptyPath(t *testing.T) {
	index, err := hnswgo.LoadIndex("", 2, "l2", uint32(10000))
	if index != nil {
		defer index.Free()
	}
	if err == nil {
		t.Fatal("No error occured when trying to load an index with a fake location.")
	} else if !os.IsNotExist(err) {
		t.Fatalf("Unexpected error received: %s", err.Error())
	}
}

func TestLoadIndexInvalidPerms(t *testing.T) {
	index, err := hnswgo.LoadIndex("/root/index.bin", 2, "l2", uint32(10000))
	if index != nil {
		defer index.Free()
	}
	if err == nil {
		t.Fatal("No error occured when trying to load an index with a fake location.")
	} else if !os.IsPermission(err) {
		t.Fatalf("Unexpected error received: %s", err.Error())
	}
}

func TestLoadIndexInvalidData(t *testing.T) {
	f, err := os.Create("/tmp/hnswgo_test.txt")
	if err != nil {
		t.Fatal(err.Error())
	}
	if _, err := f.Write([]byte("hello world")); err != nil {
		t.Fatal(err.Error())
	}

	index, err := hnswgo.LoadIndex("/tmp/hnswgo_test.txt", 2, "l2", uint32(10000))
	if index != nil {
		defer index.Free()
	}
	if err == nil {
		t.Fatal("No error occured when trying to load an index using invalid data.")
	} else if err.Error() != "Index seems to be corrupted or unsupported" {
		t.Fatalf("Unexpected error received: %s", err.Error())
	}

	// cleanup
	if err := os.Remove("/tmp/hnswgo_test.txt"); err != nil {
		t.Fatal(err.Error())
	}
}

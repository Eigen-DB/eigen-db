package tests

import (
	"errors"
	"fmt"
	"math/rand/v2"
	"os"
	"testing"

	"github.com/Eigen-DB/eigen-db/libs/faissgo/v3/faiss"
	"github.com/Eigen-DB/eigen-db/libs/faissgo/v3/index"
)

var testIndex index.Index

func generateRandomVectors(numVecs int, dim int) []float32 {
	vectors := make([]float32, dim*numVecs)
	for i := 0; i < numVecs; i++ {
		for j := 0; j < dim; j++ {
			vectors[i*dim+j] = (rand.Float32() * 2.0) - 1.0
		}
	}
	return vectors
}

func setup() error {
	idx, err := index.IndexFactory(
		128,
		"HNSW32_PQ16x8",
		faiss.MetricL2,
	)
	if err != nil {
		return errors.New("Error setting up test suite: " + err.Error())
	}
	testIndex = idx
	return nil
}

func teardown() error {
	testIndex.Free()
	if err := os.Remove(idxPersistPath); err != nil {
		return err
	}
	return nil
}

func TestMain(m *testing.M) {
	if err := setup(); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	code := m.Run()
	if err := teardown(); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	os.Exit(code)
}

func TestIndexFactory(t *testing.T) {
	idx, err := index.IndexFactory(
		128,
		"HNSW32_PQ16x8",
		faiss.MetricL2,
	)
	if err != nil {
		t.Errorf("Error creating index: %v", err)
	}
	idx.Free()
}

func TestTrain(t *testing.T) {
	vectors := generateRandomVectors(1000, 128)
	if err := testIndex.Train(vectors); err != nil {
		t.Errorf("Error training index: %v", err)
	}
}

func TestAdd(t *testing.T) {
	vectors := generateRandomVectors(1000, 128)
	if err := testIndex.Add(vectors); err != nil {
		t.Errorf("Error adding vectors: %v", err)
	}
}

func TestAddWithIds(t *testing.T) { // potential issue: overlap of ids with vectors inserted in TestAdd
	vectors := generateRandomVectors(1000, 128)
	ids := make([]int64, 1000)
	for i := 0; i < 1000; i++ {
		ids[i] = int64(i) * 7 // Add() will automatically assign IDs in increasing order so I chose a different pattern to test AddWithIds()
	}
	if err := testIndex.AddWithIds(vectors, ids); err != nil {
		t.Errorf("Error adding vectors with IDs: %v", err)
	}
}

func TestRemoveIds(t *testing.T) {
	t.Skip("Not implemented")
}

func TestSearch(t *testing.T) {
	queryVector := generateRandomVectors(1, 128)
	k := int64(5)
	ids, dists, err := testIndex.Search(queryVector, k) // vectors have already been added from TestAdd and TestAddWithIds
	if err != nil {
		t.Errorf("Error searching KNN: %v", err)
	}
	if int64(len(ids)) != k || int64(len(dists)) != k {
		t.Errorf("Expected %d NNs, got %d", k, len(ids))
	}
}

func TestReconstruct(t *testing.T) {
	_, err := testIndex.Reconstruct(0) // todo: log the reconstructed vector (curious)
	if err != nil {
		t.Errorf("Error reconstructing vector: %v", err)
	}
}

func TestIsTrained(t *testing.T) {
	if !testIndex.IsTrained() {
		t.Errorf("Index is said to not be trained when in fact it was trained in TestTrain()")
	}
}

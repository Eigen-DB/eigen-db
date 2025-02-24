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

var idxIDMap index.Index
var idxHNSW index.Index
var idxPQ index.Index

const DIM int = 128

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
	_idxIDMap, err := index.IndexFactory(
		DIM,
		"IDMap,HNSW32",
		faiss.MetricL2,
	)
	if err != nil {
		return errors.New("Error setting up test suite: " + err.Error())
	}

	_idxHNSW, err := index.IndexFactory(
		DIM,
		"HNSW32",
		faiss.MetricL2,
	)
	if err != nil {
		return errors.New("Error setting up test suite: " + err.Error())
	}

	_idxPQ, err := index.IndexFactory(
		DIM,
		"HNSW32_PQ16x8",
		faiss.MetricL2,
	)
	if err != nil {
		return errors.New("Error setting up test suite: " + err.Error())
	}

	idxIDMap = _idxIDMap
	idxHNSW = _idxHNSW
	idxPQ = _idxPQ

	return nil
}

func teardown() error {
	idxIDMap.Free()
	idxHNSW.Free()
	idxPQ.Free()
	if _, err := os.Stat(idxPersistPath); !os.IsNotExist(err) {
		if err := os.Remove(idxPersistPath); err != nil {
			return err
		}
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
	dim := 2
	idx, err := index.IndexFactory( // creating an index with smaller dimensionality to speed up training
		dim,
		"PQ2x8",
		faiss.MetricL2,
	)
	if err != nil {
		t.Errorf("Error creating index: %v", err)
	}
	defer idx.Free()
	vectors := generateRandomVectors(256, dim)
	if err := idx.Train(vectors); err != nil {
		t.Errorf("Error training index: %v", err)
	}
}

func TestAdd(t *testing.T) {
	vectors := generateRandomVectors(1000, DIM)
	if err := idxHNSW.Add(vectors); err != nil {
		t.Errorf("Error adding vectors: %v", err)
	}
}

func TestAddWithIds(t *testing.T) {
	vectors := generateRandomVectors(1000, DIM)
	ids := make([]int64, 1000)
	for i := 0; i < 1000; i++ {
		ids[i] = int64(i) * 7 // Add() will automatically assign IDs in increasing order so I chose a different pattern to test AddWithIds()
	}
	if err := idxIDMap.AddWithIds(vectors, ids); err != nil {
		t.Errorf("Error adding vectors with IDs: %v", err)
	}
}

func TestRemoveIds(t *testing.T) {
	t.Skip("Not implemented")
}

func TestSearch(t *testing.T) {
	queryVector := generateRandomVectors(1, DIM)
	k := int64(5)
	ids, dists, err := idxIDMap.Search(queryVector, k) // vectors have already been added from TestAdd and TestAddWithIds
	if err != nil {
		t.Errorf("Error searching KNN: %v", err)
	}
	if int64(len(ids)) != k || int64(len(dists)) != k {
		t.Errorf("Expected %d NNs, got %d", k, len(ids))
	}
}

func TestReconstruct(t *testing.T) {
	v := generateRandomVectors(1, DIM)
	if err := idxPQ.Add(v); err != nil {
		t.Errorf("Error adding vector: %v", err)
	}
	r, err := idxPQ.Reconstruct(0)
	t.Logf("Reconstructed vector: %v", r) // curious
	if err != nil {
		t.Errorf("Error reconstructing vector: %v", err)
	}
}

func TestIsTrained(t *testing.T) {
	if !idxPQ.IsTrained() {
		t.Errorf("Index is said to not be trained when in fact it was trained in TestTrain()")
	}
}

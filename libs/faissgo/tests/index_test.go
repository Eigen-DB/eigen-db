package tests

import (
	"fmt"
	"math/rand/v2"
	"os"
	"testing"

	"github.com/Eigen-DB/eigen-db/libs/faissgo/v0"
)

const DIM int = 128
const TEST_TMP_PATH string = "/tmp/faissgo_test"
const IDX_PERSIST_PATH string = TEST_TMP_PATH + "/index_test.bin"

// generate random vectors with values between [-1, 1)
func generateRandomVectors(numVecs int, dim int) []float32 {
	vectors := make([]float32, dim*numVecs)
	for i := 0; i < numVecs; i++ {
		for j := 0; j < dim; j++ {
			vectors[i*dim+j] = (rand.Float32() * 2.0) - 1.0
		}
	}
	return vectors
}

func getIndex(t *testing.T, indexType string, dim int) faissgo.Index {
	idx, err := faissgo.IndexFactory(
		dim,
		indexType,
		faissgo.MetricL2,
	)
	if err != nil {
		t.Errorf("Error generating test index: %s", err.Error())
	}
	return idx
}

func setup() error {
	if _, err := os.Stat(TEST_TMP_PATH); os.IsNotExist(err) {
		if err := os.MkdirAll(TEST_TMP_PATH, 0755); err != nil {
			return err
		}
	}
	return nil
}

func teardown() error {
	if _, err := os.Stat(TEST_TMP_PATH); !os.IsNotExist(err) {
		if err := os.Remove(TEST_TMP_PATH); err != nil {
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
	idx, err := faissgo.IndexFactory(
		128,
		"HNSW32_PQ16x8",
		faissgo.MetricL2,
	)
	if err != nil {
		t.Errorf("Error creating index: %v", err)
	}
	idx.Free()
}

func TestTrain(t *testing.T) {
	idx := getIndex(t, "HNSW32_PQ2x2", 2)
	defer idx.Free()
	vectors := generateRandomVectors(128, 2)
	if err := idx.Train(vectors); err != nil {
		t.Errorf("Error training index: %v", err)
	}
}

func TestAdd(t *testing.T) {
	idx := getIndex(t, "HNSW32", DIM)
	defer idx.Free()
	vectors := generateRandomVectors(1000, DIM)
	if err := idx.Add(vectors); err != nil {
		t.Errorf("Error adding vectors: %v", err)
	}
}

func TestAddWithIds(t *testing.T) {
	idx := getIndex(t, "IDMap,HNSW32", DIM)
	defer idx.Free()
	vectors := generateRandomVectors(1000, DIM)
	ids := make([]int64, 1000)
	for i := 0; i < 1000; i++ {
		ids[i] = int64(i) * 7 // just to make sure they are not sequential
	}
	if err := idx.AddWithIds(vectors, ids); err != nil {
		t.Errorf("Error adding vectors with IDs: %v", err)
	}
}

func TestRemoveIds(t *testing.T) {
	t.Skip("Not implemented")
}

func TestSearch(t *testing.T) {
	idx := getIndex(t, "HNSW32", DIM)
	defer idx.Free()

	vectors := generateRandomVectors(1000, DIM)
	if err := idx.Add(vectors); err != nil {
		t.Errorf("Error adding vectors: %v", err)
	}

	queryVector := generateRandomVectors(1, DIM)
	k := int64(5)
	ids, dists, err := idx.Search(queryVector, k)
	if err != nil {
		t.Errorf("Error searching KNN: %v", err)
	}
	if int64(len(ids)) != k || int64(len(dists)) != k {
		t.Errorf("Expected %d NNs, got %d", k, len(ids))
	}
}

func TestSearchKGreaterThanN(t *testing.T) {
	idx := getIndex(t, "HNSW32", DIM)
	defer idx.Free()

	vectors := generateRandomVectors(1000, DIM)
	if err := idx.Add(vectors); err != nil {
		t.Errorf("Error adding vectors: %v", err)
	}

	queryVector := generateRandomVectors(1, DIM)
	N := idx.NTotal()
	k := N + 1                                    // we want k > N
	ids, dists, err := idx.Search(queryVector, k) // vectors have already been added from TestAdd and TestAddWithIds
	if err != nil {
		t.Errorf("Error searching KNN: %v", err)
	}
	if int64(len(ids)) != N || int64(len(dists)) != N {
		t.Errorf("Expected %d NNs, got %d", k, len(ids))
	}
}

func TestReconstruct(t *testing.T) {
	idx := getIndex(t, "HNSW32,PQ16x2", DIM)
	defer idx.Free()

	trainVecs := generateRandomVectors(256, DIM)
	if err := idx.Train(trainVecs); err != nil {
		t.Errorf("Error training index: %v", err)
	}

	if err := idx.Add(trainVecs); err != nil {
		t.Errorf("Error adding vectors: %v", err)
	}
	r, err := idx.Reconstruct(0)
	t.Logf("Reconstructed vector: %v", r) // curious
	if err != nil {
		t.Errorf("Error reconstructing vector: %v", err)
	}
}

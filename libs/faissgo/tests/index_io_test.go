package tests

import (
	"os"
	"testing"

	"github.com/Eigen-DB/eigen-db/libs/faissgo/v3/faiss"
	"github.com/Eigen-DB/eigen-db/libs/faissgo/v3/index"
)

var idxPersistPath string = "/tmp/index.bin"

func TestWriteToDisk(t *testing.T) {
	if err := testIndex.WriteToDisk(idxPersistPath); err != nil {
		t.Fatalf("Error writing index to disk: %v", err)
	}
	if _, err := os.Stat(idxPersistPath); os.IsNotExist(err) {
		t.Fatalf("File %s does not exist after writing index to disk", idxPersistPath)
	}
}

func TestLoadFromDisk(t *testing.T) {
	idx, err := index.IndexFactory(128, "HNSW32_PQ16x8", faiss.MetricL2)
	if err != nil {
		t.Fatalf("Error creating index: %v", err)
	}
	defer idx.Free()

	if err := idx.LoadFromDisk(idxPersistPath); err != nil {
		t.Fatalf("Error loading index from disk: %v", err)
	}

	// using idx.Search() to test if the index was loaded correctly
	queryVec := generateRandomVectors(1, 128)
	k := int64(5)
	nnIds, nnDists, err := idx.Search(queryVec, k)
	if err != nil {
		t.Fatalf("Error searching KNN with index loaded from disk: %v", err)
	}
	if int64(len(nnIds)) != k || int64(len(nnDists)) != k {
		t.Fatalf("Expected %d nearest neighbors, got %d", k, len(nnIds))
	}
}

func TestLoadFromDiskInvalidIndexConfig(t *testing.T) {
	idx, err := index.IndexFactory(128, "HNSW32", faiss.MetricL2) // HNSW32 instead of HNSW32_PQ16x8
	if err != nil {
		t.Fatalf("Error creating index: %v", err)
	}
	defer idx.Free()

	if err := idx.LoadFromDisk(idxPersistPath); err == nil {
		t.Fatalf("Expected error loading index from disk with invalid index config")
	} else {
		t.Logf("(THIS ERROR IS ON PURPOSE) Error loading index from disk with invalid index config: %v", err)
	}
}

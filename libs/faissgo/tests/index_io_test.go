package tests

import (
	"os"
	"testing"
)

func TestWriteToDisk(t *testing.T) {
	idx := getIndex(t, "HNSW32", DIM)
	defer idx.Free()
	vectors := generateRandomVectors(10, DIM)
	if err := idx.Add(vectors); err != nil {
		t.Errorf("Error adding vectors: %v", err)
	}
	if err := idx.WriteToDisk(IDX_PERSIST_PATH); err != nil {
		t.Fatalf("Error writing index to disk: %v", err)
	}
	if _, err := os.Stat(IDX_PERSIST_PATH); os.IsNotExist(err) {
		t.Fatalf("File %s does not exist after writing index to disk", IDX_PERSIST_PATH)
	}
	// Clean up the file after the test
	if err := os.Remove(IDX_PERSIST_PATH); err != nil {
		t.Fatalf("Error removing file %s: %v", IDX_PERSIST_PATH, err)
	}
}

func TestLoadFromDisk(t *testing.T) {
	idx := getIndex(t, "HNSW32", DIM)
	defer idx.Free()

	// create persisted index
	vectors := generateRandomVectors(10, DIM)
	if err := idx.Add(vectors); err != nil {
		t.Errorf("Error adding vectors: %v", err)
	}
	if err := idx.WriteToDisk(IDX_PERSIST_PATH); err != nil {
		t.Fatalf("Error writing index to disk: %v", err)
	}

	// load persisted index into memory
	if err := idx.LoadFromDisk(IDX_PERSIST_PATH); err != nil {
		t.Fatalf("Error loading index from disk: %v", err)
	}

	// using idx.Search() to test if the index was loaded correctly
	queryVec := generateRandomVectors(1, DIM)
	k := int64(5)
	nnIds, nnDists, err := idx.Search(queryVec, k)
	if err != nil {
		t.Fatalf("Error searching KNN with index loaded from disk: %v", err)
	}
	if int64(len(nnIds)) != k || int64(len(nnDists)) != k {
		t.Fatalf("Expected %d nearest neighbors, got %d", k, len(nnIds))
	}

	// Clean up the file after the test
	if err := os.Remove(IDX_PERSIST_PATH); err != nil {
		t.Fatalf("Error removing file %s: %v", IDX_PERSIST_PATH, err)
	}
}

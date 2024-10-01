package vector_io

/*
import (
	"bytes"
	"eigen_db/constants"
	"encoding/gob"
	"errors"
	"os"
	"testing"

	"github.com/Eigen-DB/hnswgo/v2"
	"github.com/stretchr/testify/assert"
)

var mockVectorStore = &vectorStore{}

func cleanup(t *testing.T) {
	err := os.Remove(constants.TESTING_TMP_FILES_PATH + "/test_vector_space.vec")
	if err != nil {
		t.Errorf("Error cleaning up after a unit test: %s", err.Error())
	}
}

func generateDummySerializedData(t *testing.T, outputFilePath string, dummyStore *vectorStore) {
	buf := new(bytes.Buffer)
	encoder := gob.NewEncoder(buf)
	if err := encoder.Encode(dummyStore); err != nil {
		t.Errorf("Error generating dummy serialized data: %s", err.Error())
	}
	serializedData := buf.Bytes()

	if err := os.WriteFile(outputFilePath, serializedData, constants.DB_PERSIST_CHMOD); err != nil {
		t.Errorf("Error writing serialized data to file: %s", err.Error())
	}
}

func generateDummyVectorStore(t *testing.T) *vectorStore {
	index, err := hnswgo.New(2, 2, 1, 42, 100, "cosine")
	if err != nil {
		t.Fatalf("An error occured when creating index: %s", err.Error())
	}
	dummyStore := &vectorStore{
		StoredVectors: map[int]*Vector{
			1: {Id: 1, Embedding: []float32{1, 2}},
			2: {Id: 2, Embedding: []float32{3, 4}},
			3: {Id: 3, Embedding: []float32{5, 6}},
		},
		index: index,
	}
	return dummyStore
}

func TestPersistToDisk_success(t *testing.T) {
	defer cleanup(t)
	path := constants.TESTING_TMP_FILES_PATH + "/test_vector_space.vec"
	err := mockVectorStore.persistToDisk(path)
	if err != nil {
		t.Fatalf("An error occured when persisting vector store to disk: %s", err.Error())
	}

	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		t.Fatal("File containing persisted data does not exist.")
	}
}

func TestPersistToDisk_no_perms_for_path(t *testing.T) {
	path := "/root/test_vector_space.vec" // shouldn't have perms to write here (assuming this isn't being ran as root)
	err := mockVectorStore.persistToDisk(path)
	if err == nil {
		t.Fatalf("No error occured when trying to write to a path without the proper permissions.")
	} else if !errors.Is(err, os.ErrPermission) {
		t.Fatalf("An error OTHER than permission issues occured: %s", err.Error())
	}
}

func TestPersistToDisk_invalid_path(t *testing.T) {
	path := "/some/fake/path/test_vector_space.vec"
	err := mockVectorStore.persistToDisk(path)
	if err == nil {
		t.Fatalf("No error occured when trying to write to an invalid path.")
	} else if !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("An error OTHER than invalid path issues occured: %s", err.Error())
	}
}

func TestLoadPersistedVectors_success(t *testing.T) {
	defer cleanup(t)
	dummyStore := generateDummyVectorStore(t)
	persistPath := constants.TESTING_TMP_FILES_PATH + "/test_vector_space.vec"
	generateDummySerializedData(t, persistPath, dummyStore)
	if err := dummyStore.loadPersistedVectors(persistPath); err != nil {
		t.Fatalf("An error occured when loading persisted data into memory: %s", err.Error())
	}
}

func TestLoadPersistedVectors_invalid_path(t *testing.T) {
	dummyStore := generateDummyVectorStore(t)
	persistPath := "/some/fake/path/dummyData.vec"

	if err := dummyStore.loadPersistedVectors(persistPath); err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			t.Fatalf("An error OTHER than the file not exisiting occured: %s", err.Error())
		}
	} else {
		t.Fatalf("No error was produced when trying to load persisted vectors from invalid source.")
	}
}

func TestLoadPersistedVectors_no_perms_for_path(t *testing.T) {
	dummyStore := generateDummyVectorStore(t)
	persistPath := "/root/dummyData.vec"

	if err := dummyStore.loadPersistedVectors(persistPath); err != nil {
		if !errors.Is(err, os.ErrPermission) {
			t.Fatalf("An error OTHER than not having the right perms: %s", err.Error())
		}
	} else {
		t.Fatalf("No error was produced when trying to load persisted vectors frm a source which requires perms I do not have.")
	}
}

func TestLoadPersistedVectors_invalid_vector(t *testing.T) {
	index, err := hnswgo.New(2, 2, 1, 42, 100, "cosine")
	if err != nil {
		t.Fatalf("An error occured when creating index: %s", err.Error())
	}
	dummyStore := &vectorStore{
		StoredVectors: map[int]*Vector{
			1: {Id: 1, Embedding: []float32{1, 2}},
			2: {Id: 2, Embedding: []float32{3, 4, 3}}, // 3D vector should cause a panic since index is 2D
			3: {Id: 3, Embedding: []float32{5, 6}},
		},
		index: index,
	}
	defer cleanup(t)

	persistPath := constants.TESTING_TMP_FILES_PATH + "/test_vector_space.vec"
	generateDummySerializedData(t, persistPath, dummyStore)

	assert.Panics(t, func() {
		if err := dummyStore.loadPersistedVectors(persistPath); err != nil {
			t.Fatalf("An error occured instead of a panic")
		}
	}, "no panic occured when trying to load an invalid persisted vector")
}
*/

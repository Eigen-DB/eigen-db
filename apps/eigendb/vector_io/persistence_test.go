package vector_io

import (
	"bytes"
	"eigen_db/cfg"
	"eigen_db/constants"
	"encoding/gob"
	"errors"
	"os"
	"testing"

	"github.com/Eigen-DB/hnswgo/v2"
)

func cleanup(t *testing.T) {
	if err := os.Remove(constants.TESTING_TMP_FILES_PATH + "/test_vector_space.vec"); err != nil {
		t.Errorf("Error cleaning up after a unit test: %s", err.Error())
	}
	if err := os.Remove(constants.TESTING_TMP_FILES_PATH + "/test_index.bin"); err != nil {
		t.Errorf("Error cleaning up after a unit test: %s", err.Error())
	}
}

func createSampleConfig(t *testing.T) {
	sampleConfig := `persistence:
  timeInterval: 3s
api:
  port: 8080
  address: 0.0.0.0
hnswParams:
  dimensions: 2
  similarityMetric: cosine
  vectorSpaceSize: 10000
  M: 32
  efConstruction: 400`

	// write sample config yaml file
	if err := os.WriteFile(constants.TESTING_TMP_FILES_PATH+"/config.yml", []byte(sampleConfig), 0777); err != nil {
		t.Fatalf("Error writing sample config file: %s", err.Error())
	}
}

func generateDummySerializedData(t *testing.T, outputStorePath string, outputIndexPath string, dummyStore *vectorStore) {
	buf := new(bytes.Buffer)
	encoder := gob.NewEncoder(buf)
	if err := encoder.Encode(dummyStore); err != nil {
		t.Errorf("Error generating dummy serialized data: %s", err.Error())
	}
	serializedData := buf.Bytes()

	if err := os.WriteFile(outputStorePath, serializedData, constants.DB_PERSIST_CHMOD); err != nil {
		t.Errorf("Error writing serialized data to file: %s", err.Error())
	}

	if err := dummyStore.index.SaveToDisk(outputIndexPath); err != nil {
		t.Errorf("Error saving index to disk: %s", err.Error())
	}
}

func generateDummyVectorStore(t *testing.T) *vectorStore {
	index, err := hnswgo.New(2, 2, 1, 42, 100, "cosine")
	if err != nil {
		t.Fatalf("An error occured when creating index: %s", err.Error())
	}

	sampleVecs := [][]float32{
		{1, 2},
		{3, 4},
		{5, 6},
	}

	dummyStore := &vectorStore{
		LatestId: 0,
		index:    index,
	}

	for i, v := range sampleVecs {
		if err := dummyStore.index.InsertVector(v, uint64(i)); err != nil {
			t.Fatalf("Error inserting sample vectors in dummy index: %s", err.Error())
		}
	}

	return dummyStore
}

func TestPersistToDisk_success(t *testing.T) {
	defer cleanup(t)
	mockVectorStore := generateDummyVectorStore(t)
	spacePath := constants.TESTING_TMP_FILES_PATH + "/test_vector_space.vec"
	indexPath := constants.TESTING_TMP_FILES_PATH + "/test_index.bin"
	err := mockVectorStore.persistToDisk(spacePath, indexPath)
	if err != nil {
		t.Fatalf("An error occured when persisting vector store to disk: %s", err.Error())
	}
}

func TestPersistToDisk_no_perms_for_space_path(t *testing.T) {
	mockVectorStore := generateDummyVectorStore(t)
	spacePath := "/root/test_vector_space.vec" // shouldn't have perms to write here (assuming this isn't being ran as root)
	indexPath := constants.TESTING_TMP_FILES_PATH + "/test_index.bin"
	err := mockVectorStore.persistToDisk(spacePath, indexPath)
	if err == nil {
		t.Fatalf("No error occured when trying to write to a path without the proper permissions.")
	} else if !errors.Is(err, os.ErrPermission) {
		t.Fatalf("An error OTHER than permission issues occured: %s", err.Error())
	}
}

func TestPersistToDisk_invalid_space_path(t *testing.T) {
	mockVectorStore := generateDummyVectorStore(t)
	spacePath := "/some/fake/path/test_vector_space.vec"
	indexPath := constants.TESTING_TMP_FILES_PATH + "/test_index.bin"
	err := mockVectorStore.persistToDisk(spacePath, indexPath)
	if err == nil {
		t.Fatalf("No error occured when trying to write to an invalid path.")
	} else if !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("An error OTHER than invalid path issues occured: %s", err.Error())
	}
}

func TestPersistToDisk_no_perms_for_index_path(t *testing.T) {
	mockVectorStore := generateDummyVectorStore(t)
	spacePath := constants.TESTING_TMP_FILES_PATH + "/test_vector_space.vec"
	indexPath := "/root/index.bin"
	err := mockVectorStore.persistToDisk(spacePath, indexPath)
	if err == nil {
		t.Fatalf("No error occured when trying to write to a path without the proper permissions.")
	} else if !errors.Is(err, os.ErrPermission) {
		t.Fatalf("An error OTHER than permission issues occured: %s", err.Error())
	}
}

func TestPersistToDisk_invalid_index_path(t *testing.T) {
	mockVectorStore := generateDummyVectorStore(t)
	spacePath := constants.TESTING_TMP_FILES_PATH + "/test_vector_space.vec"
	indexPath := "/some/fake/path/index.bin"
	err := mockVectorStore.persistToDisk(spacePath, indexPath)
	if err == nil {
		t.Fatalf("No error occured when trying to write to an invalid path.")
	} else if !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("An error OTHER than invalid path issues occured: %s", err.Error())
	}
}

func TestLoadPersistedVectors_success(t *testing.T) {
	defer cleanup(t)
	createSampleConfig(t)
	dummyStore := generateDummyVectorStore(t)
	if err := cfg.SetupConfig(constants.TESTING_TMP_FILES_PATH + "/config.yml"); err != nil {
		t.Fatal(err.Error())
	}
	spacePersistPath := constants.TESTING_TMP_FILES_PATH + "/test_vector_space.vec"
	indexPersistPath := constants.TESTING_TMP_FILES_PATH + "/test_index.bin"
	generateDummySerializedData(t, spacePersistPath, indexPersistPath, dummyStore)
	if err := dummyStore.loadPersistedStore(spacePersistPath, indexPersistPath); err != nil {
		t.Fatalf("An error occured when loading persisted data into memory: %s", err.Error())
	}
}

func TestLoadPersistedVectors_invalid_store_path(t *testing.T) {
	dummyStore := generateDummyVectorStore(t)
	spacePersistPath := "/some/fake/path/dummyData.vec"
	indexPersistPath := constants.TESTING_TMP_FILES_PATH + "/test_index.bin"

	if err := dummyStore.loadPersistedStore(spacePersistPath, indexPersistPath); err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			t.Fatalf("An error OTHER than the file not exisiting occured: %s", err.Error())
		}
	} else {
		t.Fatalf("No error was produced when trying to load persisted vectors from invalid source.")
	}
}

func TestLoadPersistedVectors_no_perms_for_store_path(t *testing.T) {
	dummyStore := generateDummyVectorStore(t)
	spacePersistPath := "/root/dummyData.vec"
	indexPersistPath := constants.TESTING_TMP_FILES_PATH + "/test_index.bin"

	if err := dummyStore.loadPersistedStore(spacePersistPath, indexPersistPath); err != nil {
		if !errors.Is(err, os.ErrPermission) {
			t.Fatalf("An error OTHER than not having the right perms: %s", err.Error())
		}
	} else {
		t.Fatalf("No error was produced when trying to load persisted vectors frm a source which requires perms I do not have.")
	}
}

func TestLoadPersistedVectors_invalid_index_path(t *testing.T) {
	createSampleConfig(t)
	dummyStore := generateDummyVectorStore(t)
	spacePersistPath := constants.TESTING_TMP_FILES_PATH + "/test_vector_space.vec"
	indexPersistPath := "/some/fake/path/index.bin"
	if err := cfg.SetupConfig(constants.TESTING_TMP_FILES_PATH + "/config.yml"); err != nil {
		t.Fatal(err.Error())
	}
	generateDummySerializedData(t, spacePersistPath, constants.TESTING_TMP_FILES_PATH+"/test_index.bin", dummyStore)

	if err := dummyStore.loadPersistedStore(spacePersistPath, indexPersistPath); err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			t.Fatalf("An error OTHER than the file not exisiting occured: %s", err.Error())
		}
	} else {
		t.Fatalf("No error was produced when trying to load persisted vectors from invalid source.")
	}
}

func TestLoadPersistedVectors_no_perms_for_index_path(t *testing.T) {
	createSampleConfig(t)
	dummyStore := generateDummyVectorStore(t)
	spacePersistPath := constants.TESTING_TMP_FILES_PATH + "/test_vector_space.vec"
	indexPersistPath := "/root/index.bin"
	if err := cfg.SetupConfig(constants.TESTING_TMP_FILES_PATH + "/config.yml"); err != nil {
		t.Fatal(err.Error())
	}
	if err := dummyStore.loadPersistedStore(spacePersistPath, indexPersistPath); err != nil {
		if !errors.Is(err, os.ErrPermission) {
			t.Fatalf("An error OTHER than not having the right perms: %s", err.Error())
		}
	} else {
		t.Fatalf("No error was produced when trying to load persisted vectors frm a source which requires perms I do not have.")
	}
}

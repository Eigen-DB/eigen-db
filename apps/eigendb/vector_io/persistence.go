package vector_io

import (
	"bytes"
	"eigen_db/constants"
	"encoding/gob"
	"fmt"
	"os"
	"time"

	"eigen_db/cfg"

	"github.com/Eigen-DB/hnswgo/v2"
)

// Persists the vector store to disk.
//
// Works by storing the serialized form of the vector store as one
// file at "storePersistPath", and the index within the vector
// store as another file at "indexPersistPath".
//
// The index is persisted separately as the vector store only contains
// a pointer to the index. This means that serializing the vector store
// will only serialize the pointer to the index, and not the index itself.
//
// Returns an error if one occured.
func (store *vectorStore) persistToDisk(storePersistPath string, indexPersistPath string) error {
	// persist the actual index
	if err := store.index.SaveToDisk(indexPersistPath); err != nil {
		return err
	}

	// persist the vector store
	buf := new(bytes.Buffer)
	encoder := gob.NewEncoder(buf)
	err := encoder.Encode(store)
	if err != nil {
		return err
	}
	serializedData := buf.Bytes()

	return os.WriteFile(storePersistPath, serializedData, constants.DB_PERSIST_CHMOD)
}

// Loads a vector store persisted on disk into memory.
//
// Loads in the actual vector store from "storePersistPath" first,
// then loads in the index into the vector store from "indexPersistPath".
//
// Returns an error if one occured.
func (store *vectorStore) loadPersistedStore(storePersistPath string, indexPersistPath string) error {
	// load the store
	serializedStore, err := os.ReadFile(storePersistPath)
	if err != nil {
		return err
	}
	buf := bytes.NewBuffer(serializedStore)
	decoder := gob.NewDecoder(buf)
	err = decoder.Decode(store)
	if err != nil {
		return err
	}

	// load the index
	config := cfg.GetConfig()
	index, err := hnswgo.LoadIndex(
		indexPersistPath,
		config.GetDimensions(),
		config.GetSimilarityMetric().ToString(),
		config.GetSpaceSize(),
	)
	if err != nil {
		return err
	}

	store.index = index

	return nil
}

// Creates a Goroutine that periodically persists data to disk.
//
// The time interval at which data is persisted on disk is set
// in the "config" at persistence.timeInterval.
//
// Returns an error if one occured.
func StartPersistenceLoop(config *cfg.Config) error {
	go func() {
		for {
			err := store.persistToDisk(constants.STORE_PERSIST_PATH, constants.INDEX_PERSIST_PATH)
			if err != nil {
				fmt.Printf("Failed to persist data to disk: %s\n", err)
			}

			time.Sleep(cfg.GetConfig().GetPersistenceTimeInterval())
		}
	}()

	return nil
}

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

func (store *vectorStore) loadPersistedVectors(storePersistPath string, indexPersistPath string) error {
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
		config.GetHNSWParamsDimensions(),
		config.GetHNSWParamsSimilarityMetric(),
		config.GetHNSWParamsSpaceSize(),
	)
	if err != nil {
		return err
	}

	store.index = index

	return nil
}

func StartPersistenceLoop(config *cfg.Config) error {
	if _, err := os.Stat(constants.STORE_PERSIST_PATH); os.IsNotExist(err) {
		if err = os.MkdirAll(constants.EIGEN_DIR, constants.DB_PERSIST_CHMOD); err != nil { // perm should maybe be switched to 600 instead of 400
			return err
		}
	}

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

package vector_io

import (
	"bytes"
	"eigen_db/constants"
	"encoding/gob"
	"fmt"
	"os"
	"time"

	"eigen_db/cfg"
)

func (store *vectorStore) PersistToDisk() error {
	buf := new(bytes.Buffer)
	encoder := gob.NewEncoder(buf)
	err := encoder.Encode(store)
	if err != nil {
		return err
	}
	serializedData := buf.Bytes()

	return os.WriteFile(constants.DB_PERSIST_PATH, serializedData, constants.DB_PERSIST_CHMOD)
}

func (store *vectorStore) LoadPersistedVectors() error {
	serializedVectors, err := os.ReadFile(constants.DB_PERSIST_PATH)
	if err != nil {
		return err
	}
	buf := bytes.NewBuffer(serializedVectors)
	decoder := gob.NewDecoder(buf)
	err = decoder.Decode(store)
	if err != nil {
		return err
	}

	for id, v := range store.StoredVectors { // load deserialized stored vectors into the vector space
		err := store.vectorSpace.InsertVector(v.Components, uint32(id))
		if err != nil {
			return err
		}
	}

	return nil
}

func StartPersistenceLoop(config cfg.IConfig) error {
	if _, err := os.Stat(constants.DB_PERSIST_PATH); os.IsNotExist(err) {
		if err = os.MkdirAll(constants.EIGEN_DIR, constants.DB_PERSIST_CHMOD); err != nil { // perm should maybe be switched to 600 instead of 400
			return err
		}
	}

	go func() {
		for {
			err := vectorStoreInstance.PersistToDisk()
			//fmt.Println("Persisted vectors")
			if err != nil {
				fmt.Printf("Failed to persist data to disk: %s\n", err)
			}

			time.Sleep((&cfg.ConfigFactory{}).GetConfig().GetPersistenceTimeInterval())
		}
	}()

	return nil
}

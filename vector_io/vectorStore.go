package vector_io

import (
	"bytes"
	"eigen_db/constants"
	"encoding/gob"
	"fmt"
	"os"
	"time"
)

type vectorStore struct {
	Vectors []vector
}

var vectorStoreInstance *vectorStore // Where all vectors are stored in memory at runtime

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
	return decoder.Decode(store)
}

func (store *vectorStore) writeVector(v *vector) {
	store.Vectors = append(store.Vectors, *v)
}

func (store *vectorStore) deleteVector(vectorId uint32) { // TODO

}

func InstantiateVectorStore() {
	vectorStoreInstance = &vectorStore{}
	err := vectorStoreInstance.LoadPersistedVectors()
	if err != nil {
		fmt.Printf("Loaded empty vector store into memory -> error loading persisted vectors: %s\n", err)
		vectorStoreInstance.Vectors = make([]vector, 0)
	} else {
		fmt.Println("Loaded persisted vectors in memory.")
	}
}

func StartPersistenceLoop() error {
	if _, err := os.Stat(constants.DB_PERSIST_PATH); os.IsNotExist(err) {
		if err = os.MkdirAll(constants.EIGEN_DIR, constants.DB_PERSIST_CHMOD); err != nil {
			return err
		}
	}

	go func() {
		for {
			err := vectorStoreInstance.PersistToDisk()
			if err != nil {
				fmt.Printf("Failed to persist data to disk: %s\n", err)
			}
			time.Sleep(time.Second * 5)
		}
	}()

	return nil
}

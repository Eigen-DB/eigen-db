package vector_io

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"os"
	"time"

	"eigen_db/cfg"
	"eigen_db/constants"
	t "eigen_db/types"

	"github.com/Eigen-DB/eigen-db/libs/faissgo"
)

// the actual vector memIdx living in memory at runtime
var memIdx *memoryIndex

// Where all vectors are stored, and all operations on vectors performed.
// Stores a vector index and the ID of the vector most recently inserted.
type memoryIndex struct {
	index t.Index             // figure out how to free index from memory when program exits
	idMap map[t.VecId]t.VecId // map of embeddng IDs -> Faiss index IDs
}

func GetMemoryIndex() *memoryIndex {
	return memIdx
}

func MemoryIndexInit(dim int, similarityMetric t.SimMetric) error {
	// start with a fresh vector store
	memIdx = &memoryIndex{}
	faissMetric, err := similarityMetric.ToFaissMetricType()
	if err != nil {
		return err
	}
	index, err := faissgo.IndexFactory(
		dim,
		"HNSW32", // add PQ later
		faissMetric,
	)
	if err != nil {
		return err
	}

	memIdx.index = index
	memIdx.idMap = make(map[t.VecId]t.VecId)

	// attempt loading persisted data into the store
	if err = memIdx.loadPersistedStore(constants.STORE_PERSIST_PATH, constants.INDEX_PERSIST_PATH); err != nil {
		fmt.Printf("Loaded empty vector space into memory -> error loading persisted vectors: %s\n", err)
	} else {
		fmt.Println("Loaded persisted vectors in memory.")
	}

	return nil
}

// Gets a vector from the in-memory vector store using its ID.
//
// Returns the vector or an error if one occured.
func (idx *memoryIndex) getVector(id t.VecId) (*Embedding, error) {
	if indexId, ok := idx.idMap[id]; ok {
		embedding, err := idx.index.Reconstruct(indexId)
		if err != nil {
			return nil, err
		}
		return EmbeddingFactory(embedding, id)
	}
	return nil, fmt.Errorf("embedding with ID %d does not exist", id)
}

// Deletes a vector from the in-memory vector store using its ID.
//
// NOTE: might cause issues with the idMap, but look into it later.
//
// Returns an error if one occured.
// nolint:all
// func deleteVector(id t.VecId) error {
// 	return store.index.DeleteVector(id)
// }

// Inserts a vector from the in-memory vector store using its ID.
//
// Returns an error if one occured.
func (idx *memoryIndex) InsertVector(v *Embedding) error {
	if _, ok := idx.idMap[v.Id]; ok {
		return fmt.Errorf("embedding with ID %d already exists", v.Id)
	}
	idx.idMap[v.Id] = idx.index.NTotal()
	return idx.index.Add(v.Data)
}

// Returns the IDs of the nearest vectors or an error if one occured.
func (idx *memoryIndex) Search(queryVector *Embedding, k int64) ([]t.VecId, error) {
	nnIds, _, err := idx.index.Search(queryVector.Data, k)
	if err != nil {
		return nil, err
	}
	return nnIds, nil
}

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
func (idx *memoryIndex) persistToDisk(storePersistPath string, indexPersistPath string) error {
	// persist the actual index
	if err := idx.index.WriteToDisk(indexPersistPath); err != nil {
		return err
	}

	// persist the vector store
	buf := new(bytes.Buffer)
	encoder := gob.NewEncoder(buf)
	err := encoder.Encode(idx)
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
func (idx *memoryIndex) loadPersistedStore(storePersistPath string, indexPersistPath string) error {
	// load the store
	serializedIndex, err := os.ReadFile(storePersistPath)
	if err != nil {
		return err
	}
	buf := bytes.NewBuffer(serializedIndex)
	decoder := gob.NewDecoder(buf)
	err = decoder.Decode(idx)
	if err != nil {
		return err
	}
	return idx.index.LoadFromDisk(indexPersistPath)
}

// Creates a Goroutine that periodically persists data to disk.
//
// The time interval at which data is persisted on disk is set
// in the "config" at persistence.timeInterval.
//
// Returns an error if one occured.
func (idx *memoryIndex) StartPersistenceLoop(config *cfg.Config) error {
	go func() {
		for {
			err := memIdx.persistToDisk(constants.STORE_PERSIST_PATH, constants.INDEX_PERSIST_PATH)
			if err != nil {
				fmt.Printf("Failed to persist data to disk: %s\n", err)
			}

			time.Sleep(cfg.GetConfig().GetPersistenceTimeInterval())
		}
	}()

	return nil
}

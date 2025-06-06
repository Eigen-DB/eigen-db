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
	index    t.Index                // figure out how to free index from memory when program exits
	Metadata map[t.EmbId]t.Metadata // map of embedding IDs to metadata (implement later)
}

func GetMemoryIndex() *memoryIndex {
	return memIdx
}

func MemoryIndexInit(dim int, similarityMetric t.SimMetric) error {
	// start with a fresh vector store
	memIdx = &memoryIndex{}
	memIdx.Metadata = make(map[t.EmbId]t.Metadata)
	faissMetric, err := similarityMetric.ToFaissMetricType()
	if err != nil {
		return err
	}
	index, err := faissgo.IndexFactory(
		dim,
		"HNSW32,IDMap2", // add PQ later
		faissMetric,
	)
	if err != nil {
		return err
	}

	memIdx.index = index

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
func (idx *memoryIndex) Get(id t.EmbId) (*Embedding, error) {
	// get embedding data from index
	embedding, err := idx.index.Reconstruct(id)
	if err != nil {
		return nil, err
	}
	// get embedding metadata if it exists
	metadata := idx.Metadata[id] // we can assume the metadata exists if the embedding exists in the index
	return EmbeddingFactory(embedding, metadata, id)
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
func (idx *memoryIndex) Insert(v *Embedding) error {
	// check if the embedding already exists in the index
	embedding, err := idx.Get(v.Id) // ISSUE: even when error is stored in err, the error is printed to stdout, but doesn't cause a panic. this is a Faiss/faissgo specific issue
	if embedding != nil && err == nil {
		return fmt.Errorf("embedding with ID %d already exists", v.Id)
	}
	// insert the embedding into the index
	if err := idx.index.AddWithIds(v.Data, []t.EmbId{v.Id}); err != nil {
		return err
	}
	// store the metadata for the embedding
	idx.Metadata[v.Id] = v.Metadata
	return nil
}

// Performs nearest neighbor search on the given query vector.
// Returns a map of nearest neighbor IDs, their metadata and rank.
// Example return value:
//
//	{
//		"863": {
//			"metadata": {
//				"key1": "value1",
//				"key2": "value2"
//			},
//			"rank": 0
//		},
//		"934": {
//			"metadata": {
//				"key1": "value1",
//				"key2": "value2"
//			},
//			"rank": 1
//		},
//		...
//	}
func (idx *memoryIndex) Search(queryVector t.EmbeddingData, k int64) (map[t.EmbId]map[string]any, error) {
	nearestNeighbors := make(map[t.EmbId]map[string]any, k)

	// get the nearest neighbors IDs
	nnIds, _, err := idx.index.Search(queryVector, k)
	if err != nil {
		return nil, err
	}
	rank := 0
	for _, id := range nnIds {
		nearestNeighbors[id] = map[string]any{
			"metadata": idx.Metadata[id],
			"rank":     rank,
		}
		rank++
	}
	return nearestNeighbors, nil
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

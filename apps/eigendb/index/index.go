package index

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"math"
	"os"
	"time"

	"eigen_db/cfg"
	"eigen_db/constants"
	t "eigen_db/types"

	"github.com/Eigen-DB/eigen-db/libs/faissgo"
)

// the actual index memIdx living in memory at runtime
var memIdx *memoryIndex

const INDEX_TYPE string = "HNSW32,IDMap2"

// Where all vectors are stored, and all operations on vectors performed.
// Stores a vector index and the ID of the vector most recently inserted.
type memoryIndex struct {
	index         t.Index // figure out how to free index from memory when program exits
	Dimensions    int
	Metric        t.SimMetric
	EmbeddingsMap map[t.EmbId]t.EmbeddingData
	Metadata      map[t.EmbId]t.Metadata // map of embedding IDs to metadata
	Normalized    bool                   // whether the index is normalized or not (used for cosine similarity)
	RebuildIndex  bool
}

func GetMemoryIndex() *memoryIndex {
	return memIdx
}

func MemoryIndexInit(dim int, similarityMetric t.SimMetric) error {
	// start with a fresh vector store
	memIdx = &memoryIndex{}
	memIdx.Metadata = make(map[t.EmbId]t.Metadata)
	memIdx.EmbeddingsMap = make(map[t.EmbId]t.EmbeddingData)
	memIdx.Dimensions = dim
	memIdx.Metric = similarityMetric
	memIdx.RebuildIndex = false

	if similarityMetric.String() == "cosine" {
		memIdx.Normalized = true // cosine similarity requires normalized vectors
	} else {
		memIdx.Normalized = false
	}
	faissMetric, err := similarityMetric.ToFaissMetricType()
	if err != nil {
		return err
	}
	index, err := faissgo.IndexFactory(
		dim,
		INDEX_TYPE, // add more index types
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

// Normalizes a vector in place.
// normalize(v) = (1/|v|)*v
func (idx *memoryIndex) normalize(vector t.EmbeddingData) {
	var magnitude float32
	for i := range vector {
		magnitude += vector[i] * vector[i]
	}
	magnitude = float32(math.Sqrt(float64(magnitude)))

	for i := range vector {
		vector[i] *= 1.0 / magnitude
	}
}

func (idx *memoryIndex) rebuildIndex() error {
	fmt.Println("Rebuilding index...") // for testing
	// get all embeddings and their IDs
	var allEmbeddings t.EmbeddingData
	var allIds []t.EmbId
	for id, embedding := range idx.EmbeddingsMap {
		allEmbeddings = append(allEmbeddings, embedding...)
		allIds = append(allIds, id)
	}

	// reset the index
	faissMetric, err := idx.Metric.ToFaissMetricType()
	if err != nil {
		return err
	}
	idx.index.Free() // free old index from memory
	newIdx, err := faissgo.IndexFactory(
		idx.Dimensions,
		INDEX_TYPE,
		faissMetric,
	)
	if err != nil {
		return err
	}
	idx.index = newIdx

	// add all embeddings back into the index
	if err := idx.index.AddWithIds(allEmbeddings, allIds); err != nil {
		return err
	}
	return nil
}

// Gets a vector from the in-memory vector store using its ID.
//
// Returns the vector or an error if one occured.
func (idx *memoryIndex) Get(id t.EmbId) (*Embedding, error) {
	embedding, exists := idx.EmbeddingsMap[id]
	if !exists {
		return nil, fmt.Errorf("embedding with ID %d does not exist", id)
	}
	metadata := idx.Metadata[id] // we can assume the metadata exists if the embedding exists in the index
	return EmbeddingFactory(embedding, metadata, id)
}

// Deletes a vector from the in-memory vector store using its ID.
//
// NOTE: might cause issues with the idMap, but look into it later.
//
// Returns an error if one occured.
// nolint:all
func (idx *memoryIndex) Delete(id t.EmbId) error {
	_, exists := idx.EmbeddingsMap[id]
	if !exists {
		return fmt.Errorf("embedding with ID %d does not exist", id)
	}
	delete(idx.EmbeddingsMap, id)
	delete(idx.Metadata, id)
	idx.RebuildIndex = true
	return nil
}

// Inserts a vector from the in-memory vector store using its ID.
//
// Returns an error if one occured.
func (idx *memoryIndex) Insert(v *Embedding) error {
	// check if the embedding already exists in the index
	_, exists := idx.EmbeddingsMap[v.Id]
	if exists {
		return fmt.Errorf("embedding with ID %d already exists", v.Id)
	}

	// if the index is normalized, normalize the embedding
	if idx.Normalized {
		idx.normalize(v.Data)
	}
	// insert the embedding into the index
	if err := idx.index.AddWithIds(v.Data, []t.EmbId{v.Id}); err != nil {
		return err
	}
	idx.EmbeddingsMap[v.Id] = v.Data // keep the index and embedding map in sync
	idx.Metadata[v.Id] = v.Metadata
	return nil
}

func (idx *memoryIndex) Upsert(v *Embedding) error {
	// if the index is normalized, normalize the embedding
	if idx.Normalized {
		idx.normalize(v.Data)
	}
	idx.EmbeddingsMap[v.Id] = v.Data
	idx.Metadata[v.Id] = v.Metadata
	idx.RebuildIndex = true
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
	if idx.RebuildIndex {
		if err := idx.rebuildIndex(); err != nil {
			return nil, err
		}
	}

	nearestNeighbors := make(map[t.EmbId]map[string]any, k)
	// if the index is normalized, normalize the query vector
	if idx.Normalized {
		idx.normalize(queryVector)
	}
	// get the nearest neighbors IDs
	nnIds, _, err := idx.index.Search(queryVector, k) // maybe check if the order is correct by getting the distances as well
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
	idx.RebuildIndex = false
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

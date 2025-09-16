package index

import (
	"fmt"
	"math"

	"eigen_db/constants"
	t "eigen_db/types"

	"github.com/Eigen-DB/eigen-db/libs/faissgo"
)

// A representation of a vector index.
// Manages embeddings in memory using Faiss.
type Index struct {
	faissIndex     t.FaissIndex // figure out how to free index from memory when program exits
	Name           string
	IndexTypeFaiss string
	Dimensions     int
	Metric         t.SimMetric
	EmbeddingsMap  map[t.EmbId]t.EmbeddingData
	Metadata       map[t.EmbId]t.Metadata // map of embedding IDs to metadata
	Normalized     bool                   // whether the index is normalized or not (used for cosine similarity)
	RebuildIndex   bool
}

func IndexFactory(name string, dim int, similarityMetric t.SimMetric) (*Index, error) {
	idx := &Index{}
	idx.Name = name
	idx.Metadata = make(map[t.EmbId]t.Metadata)
	idx.EmbeddingsMap = make(map[t.EmbId]t.EmbeddingData)
	idx.IndexTypeFaiss = constants.INDEX_TYPE_FAISS
	idx.Dimensions = dim
	idx.Metric = similarityMetric
	idx.RebuildIndex = false

	if similarityMetric.String() == "cosine" {
		idx.Normalized = true // cosine similarity requires normalized vectors
	} else {
		idx.Normalized = false
	}
	faissMetric, err := similarityMetric.ToFaissMetricType()
	if err != nil {
		return nil, err
	}
	index, err := faissgo.IndexFactory(
		dim,
		constants.INDEX_TYPE_FAISS,
		faissMetric,
	)
	if err != nil {
		return nil, err
	}

	idx.faissIndex = index
	return idx, nil
}

// Normalizes a vector in place.
// normalize(v) = (1/|v|)*v
func (idx *Index) normalize(vector t.EmbeddingData) {
	var magnitude float32
	for i := range vector {
		magnitude += vector[i] * vector[i]
	}
	magnitude = float32(math.Sqrt(float64(magnitude)))

	for i := range vector {
		vector[i] *= 1.0 / magnitude
	}
}

func (idx *Index) rebuildIndex() error {
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
	idx.faissIndex.Free() // free old index from memory
	newIdx, err := faissgo.IndexFactory(
		idx.Dimensions,
		idx.IndexTypeFaiss,
		faissMetric,
	)
	if err != nil {
		return err
	}
	idx.faissIndex = newIdx

	// add all embeddings back into the index
	if err := idx.faissIndex.AddWithIds(allEmbeddings, allIds); err != nil {
		return err
	}
	return nil
}

// Gets a vector from the in-memory vector store using its ID.
//
// Returns the vector or an error if one occured.
func (idx *Index) Get(id t.EmbId) (*Embedding, error) {
	embedding, exists := idx.EmbeddingsMap[id]
	if !exists {
		return nil, fmt.Errorf("embedding with ID %d does not exist", id)
	}
	metadata := idx.Metadata[id] // we can assume the metadata exists if the embedding exists in the index
	return EmbeddingFactory(embedding, metadata, id), nil
}

// Deletes a vector from the in-memory vector store using its ID.
//
// NOTE: might cause issues with the idMap, but look into it later.
//
// Returns an error if one occured.
// nolint:all
func (idx *Index) Delete(id t.EmbId) error {
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
func (idx *Index) Insert(v *Embedding) error {
	if len(v.Data) != idx.Dimensions {
		return fmt.Errorf("embedding has %d dimensions while the index is %d-dimensional", len(v.Data), idx.Dimensions)
	}

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
	if err := idx.faissIndex.AddWithIds(v.Data, []t.EmbId{v.Id}); err != nil {
		return err
	}
	idx.EmbeddingsMap[v.Id] = v.Data // keep the index and embedding map in sync
	idx.Metadata[v.Id] = v.Metadata
	return nil
}

func (idx *Index) Upsert(v *Embedding) error {
	if len(v.Data) != idx.Dimensions {
		return fmt.Errorf("embedding has %d dimensions while the index is %d-dimensional", len(v.Data), idx.Dimensions)
	}

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
func (idx *Index) Search(queryVector t.EmbeddingData, k int64) (map[t.EmbId]map[string]any, error) {
	if len(queryVector) != idx.Dimensions {
		return nil, fmt.Errorf("query vector has %d dimensions while the index is %d-dimensional", len(queryVector), idx.Dimensions)
	}

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
	nnIds, _, err := idx.faissIndex.Search(queryVector, k) // maybe check if the order is correct by getting the distances as well
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

func (idx *Index) GetFaissIndex() t.FaissIndex {
	return idx.faissIndex
}

func (idx *Index) SetFaissIndex(faissIndex t.FaissIndex) {
	idx.faissIndex = faissIndex
}

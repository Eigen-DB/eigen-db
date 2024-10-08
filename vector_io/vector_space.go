package vector_io

import (
	"fmt"
	"time"

	"eigen_db/constants"
	t "eigen_db/types"

	"github.com/Eigen-DB/hnswgo/v2"
)

// the actual vector store living in memory at runtime
var store *vectorStore

// Where all vectors are stored, and all operations on vectors performed.
// Stores a vector index and the ID of the vector most recently inserted.
type vectorStore struct {
	index    t.Index
	LatestId t.VectorId
}

// Gets a vector from the in-memory vector store using its ID.
//
// Returns the vector or an error if one occured.
func getVector(id t.VectorId) (*Vector, error) {
	vector, err := store.index.GetVector(id)
	if err != nil {
		return nil, err
	}

	v, err := NewVector(vector)
	if err != nil {
		return nil, err
	}

	return v, nil
}

// Deletes a vector from the in-memory vector store using its ID.
//
// Returns an error if one occured.
func deleteVector(id t.VectorId) error {
	return store.index.DeleteVector(id)
}

// Inserts a vector from the in-memory vector store using its ID.
//
// Returns an error if one occured.
func InsertVector(v *Vector) error {
	err := store.index.InsertVector(v.Embedding, uint64(v.Id))
	if err != nil {
		return err
	}
	return nil
}

// we perform similarity search using the HNSW algorithm with a time complexity of O(log n)
// when performing the algorithm, we use k+1 as the resulting k-nearest neighbors will always include the query vector itself.
// therefore we simply perform the search for k+1 nearest neighbors and remove the queryVectorId from the output

// Performs similarity search using the query vector's ID and the specified k value.
//
// Similarity search is performed using the HNSW algorithm with a time complexity of O(log n).
// When performing the algorithm, we always find the k+1 nearesr neighbors as the k-nearest neighbors
// will always include the query vector itself. Therefore we simply perform the search for k+1 nearest
// neighbors and remove the query vector from the output.
//
// Returns the IDs of the nearest vectors or an error if one occured.
func SimilaritySearch(queryVectorId t.VectorId, k int) ([]t.VectorId, error) {
	queryVector, err := getVector(queryVectorId)
	if err != nil {
		return nil, err
	}

	ids, _, err := store.index.SearchKNN(queryVector.Embedding, k+1) // returns ids of resulting vectors and the vectors' distances from the query vector
	if err != nil {
		return nil, err
	}

	// might just need to pop the first or last neighbor if I can confirm that hnswgo will return the neighbors in order
	idsExcludingQuery := make([]t.VectorId, 0)
	for _, id := range ids {
		if id != queryVectorId {
			idsExcludingQuery = append(idsExcludingQuery, id)
		}
	}
	return idsExcludingQuery, nil
}

// Instantiates the in-memory vector store.
//
// Instantiates the store using the passed in dimension count, similarity metric,
// maximum vector count (spaceSize), m and efConstruction parameters.
//
// Attempts to load in a previously persisted vector store. If one does not
// exist, a fresh store is loaded into memory.
//
// Returns an error if one occured.
func InstantiateVectorStore(dim int, similarityMetric t.SimilarityMetric, spaceSize uint32, M int, efConstruction int) error {
	similarityMetric, err := t.ParseSimilarityMetric(similarityMetric)
	if err != nil {
		return err
	}

	store = &vectorStore{}
	index, err := hnswgo.New(
		dim,
		M,
		efConstruction,
		int(time.Now().Unix()),
		spaceSize,
		similarityMetric,
	)
	if err != nil {
		return err
	}

	store.index = index

	if err = store.loadPersistedStore(constants.STORE_PERSIST_PATH, constants.INDEX_PERSIST_PATH); err != nil {
		fmt.Printf("Loaded empty vector space into memory -> error loading persisted vectors: %s\n", err)
	} else {
		fmt.Println("Loaded persisted vectors in memory.")
	}
	return nil
}

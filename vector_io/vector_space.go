package vector_io

import (
	"fmt"
	"time"

	"eigen_db/constants"
	t "eigen_db/types"

	"github.com/Eigen-DB/hnswgo/v2"
)

var store *vectorStore // where all vectors are stored in memory at runtime

type vectorStore struct {
	index    t.Index
	LatestId t.VectorId
}

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

func deleteVector(id t.VectorId) error {
	return store.index.DeleteVector(id)
}

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

	if err = store.loadPersistedVectors(constants.STORE_PERSIST_PATH, constants.INDEX_PERSIST_PATH); err != nil {
		fmt.Printf("Loaded empty vector space into memory -> error loading persisted vectors: %s\n", err)
	} else {
		fmt.Println("Loaded persisted vectors in memory.")
	}
	return nil
}

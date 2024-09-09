package vector_io

import (
	"fmt"
	"time"

	t "eigen_db/types"

	"github.com/Eigen-DB/hnswgo"
)

var vectorStoreInstance *vectorStore // Where all vectors are stored in memory at runtime

type vectorStore struct {
	StoredVectors map[t.VectorId]*Vector
	vectorSpace   t.VectorSpace
	LatestId      t.VectorId
}

type VectorSearcher struct{}

func getVector(id t.VectorId) (*Vector, error) {
	vector := vectorStoreInstance.StoredVectors[id] // if id does not exist in StoredVectors, it returns a nil pointer
	if vector != nil {
		return vectorStoreInstance.StoredVectors[id], nil
	}
	return nil, fmt.Errorf("there is no vector with id %d in the vector space", id)
}

func (store *vectorStore) writeVector(v *Vector) {
	v.Id = vectorStoreInstance.LatestId + 1
	vectorStoreInstance.LatestId++
	store.vectorSpace.AddPoint(v.Components, uint32(v.Id))
	store.StoredVectors[v.Id] = v
}

func (searcher *VectorSearcher) SimilaritySearch(queryVectorId t.VectorId, k int) ([]t.VectorId, error) {
	// we perform similarity search using the HNSW algorithm with a time complexity of O(log n)
	// when performing the algorithm, we use k+1 as the resulting k-nearest neighbors will always include the query vector itself.
	// therefore we simply perform the search for k+1 nearest neighbors and remove the queryVectorId from the output

	queryVector, err := getVector(queryVectorId)
	if err != nil {
		return nil, err
	}
	ids, _ := vectorStoreInstance.vectorSpace.SearchKNN(queryVector.Components, k+1) // returns ids of resulting vectors and the vectors' distances from the query vector

	idsExcludingQuery := make([]t.VectorId, 0)
	for _, id := range ids {
		if int(id) != queryVectorId {
			idsExcludingQuery = append(idsExcludingQuery, int(id))
		}
	}
	return idsExcludingQuery, nil
}

func instantiateVectorStore(dim int, similarityMetric t.SimilarityMetric, spaceSize uint32, M int, efConstruction int) {
	vectorStoreInstance = &vectorStore{}
	vectorStoreInstance.vectorSpace = hnswgo.New(
		dim,
		M,
		efConstruction,
		int(time.Now().Unix()),
		spaceSize,
		similarityMetric,
	)
	vectorStoreInstance.StoredVectors = make(map[int]*Vector)

	err := vectorStoreInstance.LoadPersistedVectors()
	if err != nil {
		fmt.Printf("Loaded empty vector space into memory -> error loading persisted vectors: %s\n", err)
	} else {
		fmt.Println("Loaded persisted vectors in memory.")
	}
}

package vector_io

import (
	"fmt"
	"time"

	t "eigen_db/types"

	"github.com/evan176/hnswgo"
)

var vectorStoreInstance *vectorStore // Where all vectors are stored in memory at runtime

type vectorStore struct {
	StoredVectors map[t.VectorId]*Vector
	vectorSpace   t.VectorSpace
	LatestId      t.VectorId
}

func (store *vectorStore) writeVector(v *Vector) {
	v.Id = vectorStoreInstance.LatestId + 1
	vectorStoreInstance.LatestId++
	store.vectorSpace.AddPoint(v.Components, v.Id)
	store.StoredVectors[v.Id] = v
}

func InstantiateVectorStore(dim uint32, similarityMetric t.SimilarityMetric, spaceSize uint32, M uint32, efConstruction uint32) {
	vectorStoreInstance = &vectorStore{}
	vectorStoreInstance.vectorSpace = hnswgo.New(int(dim), int(M), int(efConstruction), int(time.Now().Unix()), spaceSize, similarityMetric)
	vectorStoreInstance.StoredVectors = make(map[uint32]*Vector)

	err := vectorStoreInstance.LoadPersistedVectors()
	if err != nil {
		fmt.Printf("Loaded empty vector space into memory -> error loading persisted vectors: %s\n", err)
	} else {
		fmt.Println("Loaded persisted vectors in memory.")
	}
}

func SimilaritySearch(queryVectorId t.VectorId, k uint32) []t.VectorId {
	// we perform similarity search using the HNSW algorithm with a time complexity of O(log n)
	// when performing the algorithm, we use k+1 as the resulting k-nearest neighbors will always include the query vector itself.
	// therefore we simply perform the search for k+1 nearest neighbors and remove the queryVectorId from the output

	queryVector := GetVector(queryVectorId)
	ids, _ := vectorStoreInstance.vectorSpace.SearchKNN(queryVector.Components, int(k)+1) // returns ids of resulting vectors and the vectors' distances from the query vector

	idsExcludingQuery := make([]t.VectorId, 0)
	for _, id := range ids {
		if id != queryVectorId {
			idsExcludingQuery = append(idsExcludingQuery, id)
		}
	}
	return idsExcludingQuery
}

func GetVector(id t.VectorId) *Vector {
	return vectorStoreInstance.StoredVectors[id]
}

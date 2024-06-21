package vector_io

import (
	"bytes"
	"eigen_db/constants"
	"encoding/gob"
	"fmt"
	"os"
	"time"

	"github.com/evan176/hnswgo"
)

type vectorStore struct {
	StoredVectors map[VectorId]*Vector
	vectorSpace   VectorSpace
	LatestId      VectorId
	Config        VectorSpaceConfig
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
	err = decoder.Decode(store)
	if err != nil {
		return err
	}

	for id, v := range store.StoredVectors { // load deserialized stored vectors into the vector space
		store.vectorSpace.AddPoint(v.Components, id)
	}

	return nil
}

func (store *vectorStore) writeVector(v *Vector) {
	store.vectorSpace.AddPoint(v.Components, v.Id)
	store.StoredVectors[v.Id] = v
}

func (space *vectorStore) deleteVector(vector *Vector) { // TODO
	// swap element at 'vector' and last element of the slice, and shrink the slice by 1, deleting the last element.
	// PROBLEM: this will make the vector IDs no longer ordered which will mess up the code that generates a vector ID when creating a new vector.
	// solution: keep the latest vector ID as a field in the vectorStore struct.
}

func InstantiateVectorStore(usePersistence bool, dim uint32, similarityMetric SimilarityMetric, spaceSize uint32, M uint32, efConstruction uint32) {
	vectorStoreInstance = &vectorStore{}
	vectorStoreInstance.vectorSpace = hnswgo.New(int(dim), int(M), int(efConstruction), int(time.Now().Unix()), spaceSize, similarityMetric)
	vectorStoreInstance.StoredVectors = make(map[uint32]*Vector)

	if usePersistence {
		err := vectorStoreInstance.LoadPersistedVectors()
		if err != nil {
			fmt.Printf("Loaded empty vector space into memory -> error loading persisted vectors: %s\n", err)
		} else {
			fmt.Println("Loaded persisted vectors in memory.")
		}
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

func SimilaritySearch(queryVectorId VectorId, k uint32) []VectorId {
	// we perform similarity search using the HNSW algorithm with a time complexity of O(log n)
	// when performing the algorithm, we use k+1 as the resulting k-nearest neighbors will always include the query vector itself.
	// therefore we simply perform the search for k+1 nearest neighbors and remove the queryVectorId from the output

	queryVector := GetVector(queryVectorId)
	ids, _ := vectorStoreInstance.vectorSpace.SearchKNN(queryVector.Components, int(k)+1) // returns ids of resulting vectors and the vectors' distances from the query vector

	idsExcludingQuery := make([]VectorId, 0)
	for _, id := range ids {
		if id != queryVectorId {
			idsExcludingQuery = append(idsExcludingQuery, id)
		}
	}
	return idsExcludingQuery
}

// GETTERS

func GetVectorStoreConfig() VectorSpaceConfig {
	return vectorStoreInstance.Config
}

func GetVector(id VectorId) *Vector {
	return vectorStoreInstance.StoredVectors[id]
}

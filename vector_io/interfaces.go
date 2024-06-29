package vector_io

import t "eigen_db/types"

type IVectorFactory interface {
	NewVector(t.VectorComponents) (IVector, error)
}

type IVector interface {
	Insert()
}

type IVectorSearcher interface {
	SimilaritySearch(t.VectorId, uint32) ([]t.VectorId, error)
}

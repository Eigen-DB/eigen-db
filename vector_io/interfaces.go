package vector_io

import t "eigen_db/types"

type IVectorFactory interface {
	NewVector([]t.VectorComponent) (IVector, error)
}

type IVector interface {
	Insert() error
}

type IVectorSearcher interface {
	SimilaritySearch(t.VectorId, int) ([]t.VectorId, error)
}

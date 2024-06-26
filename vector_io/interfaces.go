package vector_io

import t "eigen_db/types"

type IVectorFactory interface {
	NewVector(t.VectorComponents) IVector
}

type IVector interface {
	Insert()
}
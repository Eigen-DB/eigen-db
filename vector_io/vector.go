package vector_io

import t "eigen_db/types"

type Vector struct {
	Id         t.VectorId         `json:"id"`
	Components t.VectorComponents `json:"components"`
}

func (v *Vector) Insert() {
	vectorStoreInstance.writeVector(v)
}

func NewVector(components t.VectorComponents) *Vector {
	v := &Vector{}
	v.Components = components
	return v
}

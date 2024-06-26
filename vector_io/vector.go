package vector_io

import t "eigen_db/types"

type Vector struct {
	Id         t.VectorId         `json:"id"`
	Components t.VectorComponents `json:"components"`
}

func (v *Vector) Insert() {
	vectorStoreInstance.writeVector(v)
}

type VectorFactory struct{}

func (factory *VectorFactory) NewVector(components t.VectorComponents) IVector {
	v := &Vector{}
	v.Components = components
	return v
}

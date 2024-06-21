package vector_io

type Vector struct {
	Id         VectorId         `json:"id"`
	Components VectorComponents `json:"components"`
}

func (v *Vector) Insert() {
	vectorStoreInstance.writeVector(v)
}

func NewVector(components VectorComponents) *Vector {
	v := &Vector{}
	v.Components = components
	return v
}

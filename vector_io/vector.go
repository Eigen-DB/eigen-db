package vector_io

type Vector struct {
	Id         VectorId         `json:"id"`
	Components VectorComponents `json:"components"`
}

func (v *Vector) Insert() {
	vectorStoreInstance.writeVector(v)
}

func (v *Vector) Delete() {
	vectorStoreInstance.deleteVector(v)
}

func NewVector(components VectorComponents) *Vector {
	v := &Vector{}
	id := vectorStoreInstance.LatestId + 1
	vectorStoreInstance.LatestId = id
	v.Id = id
	v.Components = components
	return v
}

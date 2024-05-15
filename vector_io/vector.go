package vector_io

type vector struct {
	Id         uint32    `json:"id"`
	Components []float64 `json:"components"`
}

func (v *vector) Insert() {
	vectorStoreInstance.writeVector(v)
}

func (v *vector) Delete() {
	vectorStoreInstance.deleteVector(v.Id)
}

func (v *vector) Search(k uint32) []*vector {
	return make([]*vector, 0)
}

func NewVector(components []float64) *vector {
	v := &vector{}
	id := uint32(0)

	if len(vectorStoreInstance.Vectors) > 0 {
		lastVectorID := vectorStoreInstance.Vectors[len(vectorStoreInstance.Vectors)-1].Id
		id = lastVectorID + 1
	}

	v.Id = id
	v.Components = components
	return v
}

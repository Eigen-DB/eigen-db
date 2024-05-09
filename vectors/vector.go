package vectors

type Vector struct {
	id         uint32
	components []float32
}

func (v *Vector) Insert() error {
	return writeVector(v)
}

func (v *Vector) Delete() error {
	return deleteVector(v.id)
}

func (v *Vector) Search(k uint32) []*Vector {
	return make([]*Vector, 0)
}

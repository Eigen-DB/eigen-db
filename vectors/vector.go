package vectors

type vector struct {
	id         uint32
	components []float32
}

func (v *vector) insert() error {
	return writeVectorToMemory(v)
}

func (v *vector) search(k uint32) []*vector {
	return make([]*vector, 0)
}

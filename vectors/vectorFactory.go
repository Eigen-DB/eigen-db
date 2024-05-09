package vectors

func NewVector(components []float32) *Vector {
	v := new(Vector)
	id := uint32(1) // generate the next ID (increment)
	v.id = id
	v.components = components
	return v
}

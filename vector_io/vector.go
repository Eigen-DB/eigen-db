package vector_io

import (
	"eigen_db/cfg"
	t "eigen_db/types"
	"fmt"
)

// A representation of a vector.
//
// Each vector has an ID and an embedding.
type Vector struct {
	Id        t.VecId     `json:"id" binding:"required"`
	Embedding t.Embedding `json:"embedding" binding:"required"`
}

// Creates a new vector with the specified embedding.
//
// Returns a pointer to the new Vector, or an error if one occured.
func NewVector(embedding t.Embedding, id t.VecId) (*Vector, error) {
	dimensions := cfg.GetConfig().GetDimensions()
	if len(embedding) == dimensions {
		v := &Vector{}
		v.Embedding = embedding
		v.Id = id
		return v, nil
	}
	return nil, fmt.Errorf("provided a %d-dimensional vector while the vector space is %d-dimensional", len(embedding), dimensions)
}

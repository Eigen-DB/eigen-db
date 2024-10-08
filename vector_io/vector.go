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
	Id        t.VectorId          `json:"id"`
	Embedding []t.VectorComponent `json:"components"`
}

// Creates a new vector with the specified embedding.
//
// Returns a pointer to the new Vector, or an error if one occured.
func NewVector(embedding []t.VectorComponent) (*Vector, error) {
	dimensions := cfg.GetConfig().GetHNSWParamsDimensions()
	if len(embedding) == dimensions {
		v := &Vector{}
		v.Embedding = embedding
		v.Id = store.LatestId + 1
		store.LatestId++
		return v, nil
	}
	return nil, fmt.Errorf("provided a %d-dimensional vector while the vector space is %d-dimensional", len(embedding), dimensions)
}

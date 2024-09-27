package vector_io

import (
	"eigen_db/cfg"
	t "eigen_db/types"
	"fmt"
)

type Vector struct {
	Id        t.VectorId          `json:"id"`
	Embedding []t.VectorComponent `json:"components"`
}

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

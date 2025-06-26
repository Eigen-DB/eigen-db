package vector_io

import (
	"eigen_db/cfg"
	t "eigen_db/types"
	"fmt"
)

// A representation of a vector.
//
// Each vector has an ID and an embedding.
type Embedding struct {
	Id       t.EmbId         `json:"id" binding:"required"`
	Data     t.EmbeddingData `json:"data" binding:"required"`
	Metadata t.Metadata      `json:"metadata" binding:"required"`
}

// Creates a new vector with the specified embedding.
//
// Returns a pointer to the new Vector, or an error if one occured.
func EmbeddingFactory(data t.EmbeddingData, metadata t.Metadata, id t.EmbId) (*Embedding, error) {
	dimensions := cfg.GetConfig().GetDimensions()
	if len(data) == dimensions {
		return &Embedding{
			Id:       id,
			Data:     data,
			Metadata: metadata,
		}, nil
	}
	return nil, fmt.Errorf("provided a %d-dimensional embedding while the index is %d-dimensional", len(data), dimensions)
}

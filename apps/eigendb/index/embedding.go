package index

import (
	t "eigen_db/types"
)

// A representation of an embedding.
//
// Each embedding has an ID and an embedding.
type Embedding struct {
	Id       t.EmbId         `json:"id" binding:"required"`
	Data     t.EmbeddingData `json:"data" binding:"required"`
	Metadata t.Metadata      `json:"metadata" binding:"required"`
}

// Creates a new embedding with the specified embedding.
//
// Returns a pointer to the new Embedding
func EmbeddingFactory(data t.EmbeddingData, metadata t.Metadata, id t.EmbId) *Embedding {
	return &Embedding{
		Id:       id,
		Data:     data,
		Metadata: metadata,
	}
}

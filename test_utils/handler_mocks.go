package test_utils

import (
	"eigen_db/types"
	"eigen_db/vector_io"
	"fmt"
)

var NewVectorInvocations int = 0
var InsertInvocations int = 0
var SimilaritySearchInvocations int = 0

type MockVectorFactory struct {
	Dimensions int
}

type MockVector struct{}

type MockVectorSearcher struct{}

func (factory *MockVectorFactory) NewVector(components []types.VectorComponent) (vector_io.IVector, error) {
	NewVectorInvocations++
	if len(components) == factory.Dimensions {
		return &MockVector{}, nil
	}
	return nil, fmt.Errorf("provided a %d-dimensional vector while the vector space is %d-dimensional", len(components), factory.Dimensions)
}

func (vector *MockVector) Insert() {
	InsertInvocations++
}

func (searcher *MockVectorSearcher) SimilaritySearch(queryVectorId types.VectorId, k int) ([]types.VectorId, error) {
	SimilaritySearchInvocations++
	return []types.VectorId{1, 2, 3}, nil
}

func Cleanup() {
	NewVectorInvocations = 0
	InsertInvocations = 0
	SimilaritySearchInvocations = 0
}

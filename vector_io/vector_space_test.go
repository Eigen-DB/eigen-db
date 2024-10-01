package vector_io

import "testing"

func TestInstatiateVectorStore_success(t *testing.T) {
	if err := InstantiateVectorStore(
		2,
		"euclidean",
		100,
		2,
		400,
	); err != nil {
		t.Fatalf("Error instantiating vector store: %s", err.Error())
	}
}

func TestInstatiateVectorStore_invalid_dim(t *testing.T) {
	if err := InstantiateVectorStore(
		-2,
		"euclidean",
		100,
		2,
		400,
	); err == nil {
		t.Fatalf("Invalid dimensions error not produced")
	} else if err.Error() != "dimension must be >= 1" {
		t.Fatalf("Invalid dimensions error not produced: %s", err.Error())
	}
}

func TestInstatiateVectorStore_invalid_sim_metric(t *testing.T) {
	if err := InstantiateVectorStore(
		2,
		"x",
		100,
		2,
		400,
	); err == nil {
		t.Fatalf("Invalid similarity metric error not produced")
	} else if err.Error() != "invalid similarity metric" {
		t.Fatalf("Invalid similarity metric error not produced: %s", err.Error())
	}
}

func TestInstatiateVectorStore_invalid_m(t *testing.T) {
	if err := InstantiateVectorStore(
		2,
		"euclidean",
		100,
		-1,
		400,
	); err == nil {
		t.Fatalf("Invalid m error not produced")
	} else if err.Error() != "m must be >= 2" {
		t.Fatalf("Invalid m error not produced: %s", err.Error())
	}
}

func TestInstatiateVectorStore_invalid_ef_construction(t *testing.T) {
	if err := InstantiateVectorStore(
		2,
		"euclidean",
		100,
		2,
		-1,
	); err == nil {
		t.Fatalf("Invalid efConstruction error not produced")
	} else if err.Error() != "efConstruction must be >= 0" {
		t.Fatalf("Invalid efConstruction error not produced: %s", err.Error())
	}
}

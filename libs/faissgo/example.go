package main

import (
	"fmt"
	"math"
	"math/rand/v2"

	"github.com/Eigen-DB/eigen-db/libs/faissgo/v3/faiss"
	"github.com/Eigen-DB/eigen-db/libs/faissgo/v3/index"
)

var dim int = 128
var simMetric faiss.MetricType = faiss.MetricL2
var nBits int = 4
var bufferIdx index.Index
var mainIdx index.Index

func addVectors(vectors []float32) error { // this is a good solution TBH
	if mainIdx.IsTrained() {
		fmt.Println("Main index trained!")
		return mainIdx.Add(vectors)
	} else {
		fmt.Println("Main index NOT trained...")
		if err := bufferIdx.Add(vectors); err != nil {
			return err
		} //bufVecs = append(bufVecs, vectors...)
		n := bufferIdx.NTotal() //n := len(bufVecs) / dim

		if n >= int64(math.Pow(2, float64(nBits))) {
			fmt.Println("Swapping! Training main index...")
			defer bufferIdx.Free() // freeing buffer index from memory

			// get vectors from bufferIdx
			bufVecs, err := bufferIdx.ReconstructN(0, n)
			if err != nil {
				return err
			}

			// train mainIdx using vectors from bufferIdx
			if err := mainIdx.Train(bufVecs); err != nil {
				return err
			}
			fmt.Println("Training done.")

			// add vectors to mainIdx
			if err := mainIdx.Add(bufVecs); err != nil {
				return err
			}
		}
	}
	return nil
}

func searchANN(queryVec []float32, k int64) ([]int64, []float32, error) {
	if mainIdx.IsTrained() {
		return mainIdx.Search(queryVec, k)
	}
	return bufferIdx.Search(queryVec, k)
}

func getRandVecs(n int) []float32 {
	vectors := make([]float32, n*dim)
	for i := range n {
		for j := range dim {
			vectors[i*dim+j] = (rand.Float32() * 2.0) - 1.0 // -1.0 <= v[j] < 1.0
		}
	}
	return vectors
}

func main() {
	idx1, err := index.IndexFactory(dim, fmt.Sprintf("HNSW32_PQ16x%d", nBits), simMetric)
	if err != nil {
		panic(err)
	}
	mainIdx = idx1
	defer mainIdx.Free()

	idx2, err := index.IndexFactory(dim, "HNSW32", simMetric)
	if err != nil {
		panic(err)
	}
	bufferIdx = idx2
	defer func() {
		if !mainIdx.IsTrained() {
			bufferIdx.Free()
		}
	}()

	// first batch of vectors
	v := getRandVecs(10)
	if err := addVectors(v); err != nil {
		panic(err)
	}

	fmt.Println("n =", bufferIdx.NTotal())

	nnLabels, nnDists, err := searchANN(getRandVecs(1), 5)
	if err != nil {
		panic(err)
	}
	fmt.Println("Nearest neighbors labels: ", nnLabels)
	fmt.Println("Nearest neighbors distances: ", nnDists)

	// second batch of vectors
	v = getRandVecs(20)
	if err := addVectors(v); err != nil {
		panic(err)
	}

	nnLabels, nnDists, err = searchANN(getRandVecs(1), 5)
	if err != nil {
		panic(err)
	}
	fmt.Println("Nearest neighbors labels: ", nnLabels)
	fmt.Println("Nearest neighbors distances: ", nnDists)

	// third batch of vectors
	v = getRandVecs(100)
	if err := addVectors(v); err != nil {
		panic(err)
	}

	nnLabels, nnDists, err = searchANN(getRandVecs(1), 5)
	if err != nil {
		panic(err)
	}
	fmt.Println("Nearest neighbors labels: ", nnLabels)
	fmt.Println("Nearest neighbors distances: ", nnDists)

	// third batch of vectors
	v = getRandVecs(100)
	if err := addVectors(v); err != nil {
		panic(err)
	}

	// third batch of vectors
	v = getRandVecs(100)
	if err := addVectors(v); err != nil {
		panic(err)
	}
}

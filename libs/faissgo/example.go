package main

import (
	"fmt"
	"math/rand/v2"

	"github.com/Eigen-DB/eigen-db/libs/faissgo/v3/faiss"
	"github.com/Eigen-DB/eigen-db/libs/faissgo/v3/index"
)

func main() {
	dim := 120 // 120 % 12 = 0 ** ISSUE: dim must be divisable by 12 (m=12 for PQ12) **
	nVec := 1000000
	idx, err := index.IndexFactory(dim, "HNSW32_PQ12", faiss.MetricL2)
	if err != nil {
		panic(err)
	}
	defer idx.Free()

	vectors := make([]float32, nVec*dim)
	for i := 0; i < nVec; i++ {
		for j := 0; j < dim; j++ {
			vectors[i*dim+j] = (rand.Float32() * 2.0) - 1.0 // -1.0 <= v[j] < 1.0
		}
	}

	//fmt.Printf("%v\n", vectors)

	//fmt.Println("Training index...")
	//if err := idx.Train(vectors); err != nil { // trained on 1000 random vectors
	//	panic(err)
	//}
	//fmt.Println("Done.")
	//
	//if err := idx.WriteToDisk("./index.bin"); err != nil {
	//	panic(err)
	//}

	fmt.Println("Loading index...")
	if err := idx.LoadFromDisk("./index.bin"); err != nil {
		panic(err)
	}
	fmt.Println("Done.")
	fmt.Printf("IsTrained: %v\n", idx.IsTrained())

	// insert sample vectors
	fmt.Println("Adding vectors...")
	if err := idx.Add(vectors); err != nil { // adding vectors with IDs is not implemented for the HNSW32_PQ12 index
		panic(err)
	}
	fmt.Println("Done.")

	if err := idx.WriteToDisk("./index-1-mil.bin"); err != nil {
		panic(err)
	}

	fmt.Println("Searching KNN...")
	k := int64(3)
	ids, dists, err := idx.Search(vectors[:dim*1], k)
	if err != nil {
		panic(err)
	}
	fmt.Println("Done.")

	fmt.Printf("NN Ids: %v\n", ids)
	fmt.Printf("NN Dists: %v\n", dists)
}

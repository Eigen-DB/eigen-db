package main

import (
	"fmt"
	"math"
	"math/rand/v2"
	"sync"

	"github.com/Eigen-DB/eigen-db/libs/faissgo/v3/faiss"
	"github.com/Eigen-DB/eigen-db/libs/faissgo/v3/index"
)

var dim int = 128
var simMetric faiss.MetricType = faiss.MetricL2
var nBits int = 4
var bufferIdx index.Index
var mainIdx index.Index
var swapInitiated bool = false
var swapDone bool = false
var swapN int64

func swap(wg *sync.WaitGroup) {
	//defer bufferIdx.Free() // freeing buffer index from memory
	defer wg.Done()
	n := bufferIdx.NTotal()

	// get vectors from bufferIdx
	bufVecs, err := bufferIdx.ReconstructN(0, n)
	if err != nil {
		panic(err)
	}

	// train mainIdx using vectors from bufferIdx
	if err := mainIdx.Train(bufVecs); err != nil {
		panic(err)
	}
	fmt.Println("Training done.")

	// add vectors to mainIdx
	if err := mainIdx.Add(bufVecs); err != nil {
		panic(err)
	}
	fmt.Println("[SWAP] Done. n =", mainIdx.NTotal())

	swapDone = true
}

func moveRemainingVectors(wg *sync.WaitGroup) {
	wg.Wait() // wait for swap to finish
	if bufferIdx.NTotal() == swapN {
		fmt.Println("[REMAIN] No vectors to add to main index.")
		return
	}
	fmt.Println("[REMAIN] Adding remaining vectors to main index...")
	remainderBufVecs, err := bufferIdx.ReconstructN(swapN, bufferIdx.NTotal()-swapN)
	if err != nil {
		panic(err) //fmt.Println("Error reconstructing remaining vectors from buffer index:", err)
	}
	if err := mainIdx.Add(remainderBufVecs); err != nil {
		panic(err) //fmt.Println("Error adding vectors to main index:", err)
	}
	fmt.Println("[REMAIN] Done.")
}

// once mainIdx is trained, operations swap over to mainIdx
// while operations are performed on mainIdx, vectors inserted in bufferIdx during the swapping process are moved to mainIdx
// this means that once the swap is done, there could be a small interval of time where ANN searches won't show some of the remainder vectors in bufferIdx that haven't been moved to mainIdx yet
func addVectors(vectors []float32, wg *sync.WaitGroup) error {
	if swapDone { // need a better flag for swapping to mainIdx than .IsTrained()
		fmt.Println("[MAIN] MAIN INDEX TRAINED!")
		return mainIdx.Add(vectors)
	} else { // using buffer index
		fmt.Println("[MAIN] Main index NOT trained...")
		if err := bufferIdx.Add(vectors); err != nil {
			return err
		}
		n := bufferIdx.NTotal()
		if n >= int64(math.Pow(2, float64(nBits))) && !swapInitiated { // ready for swap
			fmt.Println("[MAIN] Swapping! Training main index...")
			swapN = n
			swapInitiated = true
			wg.Add(1)
			go swap(wg)                 // handle returning errors
			go moveRemainingVectors(wg) // move remaning vectors to main index
		}
	}
	return nil
}

func searchANN(queryVec []float32, k int64) ([]int64, []float32, error) {
	if swapDone {
		fmt.Println("[MAIN] Searching on main index...")
		return mainIdx.Search(queryVec, k)
	}
	fmt.Println("[MAIN] Searching on buffer index...")
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

	var wg sync.WaitGroup

	// first batch of vectors
	v := getRandVecs(10)
	if err := addVectors(v, &wg); err != nil {
		panic(err)
	}

	fmt.Println("n =", bufferIdx.NTotal())

	nnLabels, nnDists, err := searchANN(getRandVecs(1), 5)
	if err != nil {
		panic(err)
	}
	fmt.Println("[MAIN] Nearest neighbors labels: ", nnLabels)
	fmt.Println("[MAIN] Nearest neighbors distances: ", nnDists)

	// second batch of vectors
	v = getRandVecs(20)
	if err := addVectors(v, &wg); err != nil {
		panic(err)
	}

	fmt.Println("[MAIN] Searching ANN")
	nnLabels, nnDists, err = searchANN(getRandVecs(1), 5)
	if err != nil {
		panic(err)
	}
	fmt.Println("[MAIN] Nearest neighbors labels: ", nnLabels)
	fmt.Println("[MAIN] Nearest neighbors distances: ", nnDists)

	// third batch of vectors
	v = getRandVecs(100)
	if err := addVectors(v, &wg); err != nil {
		panic(err)
	}

	nnLabels, nnDists, err = searchANN(getRandVecs(1), 5)
	if err != nil {
		panic(err)
	}
	fmt.Println("[MAIN] Nearest neighbors labels: ", nnLabels)
	fmt.Println("[MAIN] Nearest neighbors distances: ", nnDists)

	//time.Sleep(10 * time.Second) // wait for swap to finish

	fmt.Println(mainIdx.NTotal()) // should be 130

	// third batch of vectors
	v = getRandVecs(100)
	if err := addVectors(v, &wg); err != nil {
		panic(err)
	}

	fmt.Println(mainIdx.NTotal())

	// third batch of vectors
	v = getRandVecs(100)
	if err := addVectors(v, &wg); err != nil {
		panic(err)
	}

	fmt.Println(mainIdx.NTotal())

	// // third batch of vectors
	// v = getRandVecs(100)
	// if err := addVectors(v, &wg); err != nil {
	// 	panic(err)
	// }
}

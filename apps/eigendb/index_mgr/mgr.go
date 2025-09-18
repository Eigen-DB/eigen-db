package index_mgr

import (
	"bytes"
	"eigen_db/cfg"
	"eigen_db/constants"
	"eigen_db/index"
	"eigen_db/types"
	"encoding/gob"
	"fmt"
	"os"
	"path"
	"regexp"
	"sync"
	"time"

	"github.com/Eigen-DB/eigen-db/libs/faissgo"
)

type IndexMgr struct {
	indexes map[string]*index.Index
}

var memIdxMgr *IndexMgr

func IndexMgrInit(wg *sync.WaitGroup, e2eTestMode bool) error {
	// instantiate the singleton index manager
	memIdxMgr = &IndexMgr{
		indexes: make(map[string]*index.Index),
	}
	if e2eTestMode { // skipping the persistence loop in E2E test mode
		return nil
	}
	// start persistence loop
	return memIdxMgr.startPersistenceLoop(wg)
}

func GetIndexMgr() *IndexMgr {
	return memIdxMgr
}

// returns list of relative paths of persisted indexes (e.g. 'myindex/'), their names, and an error if one occurs
//
// the two slices should be parallel (i.e. indexes[0] corresponds to indexNames[0], etc...)
func listPersistedIndexes() ([]string, []string, error) { // NOTE: returns RELATIVE paths of indexes
	var indexes []string
	var indexNames []string
	entries, err := os.ReadDir(constants.EIGEN_DIR)
	if err != nil {
		return nil, nil, err
	}
	for _, entry := range entries {
		if entry.IsDir() {
			indexes = append(indexes, entry.Name())
			indexNames = append(indexNames, path.Base(entry.Name()))
		}
	}
	return indexes, indexNames, nil
}

func (mgr *IndexMgr) GetIndex(name string) (*index.Index, error) {
	idx, exists := mgr.indexes[name]
	if !exists {
		return nil, fmt.Errorf("index '%s' does not exist", name)
	}
	return idx, nil
}

func (mgr *IndexMgr) CreateIndex(name string, dim int, metric types.SimMetric) error {
	valid, err := regexp.Match(constants.VALID_INDEX_NAME_REGEX, []byte(name))
	if err != nil {
		return err
	}
	if !valid {
		return fmt.Errorf("index name '%s' is invalid - must be between 3-32 characters long and only contain lowercase letters, numbers, and/or dashes", name)
	}

	_, exists := mgr.indexes[name]
	if exists {
		return fmt.Errorf("index '%s' already exists", name)
	}
	idx, err := index.IndexFactory(name, dim, metric)
	if err != nil {
		return err
	}
	mgr.indexes[name] = idx
	if err := os.MkdirAll(path.Join(constants.EIGEN_DIR, name), constants.DB_PERSIST_CHMOD); err != nil {
		return err
	}
	return nil
}

func (mgr *IndexMgr) DeleteIndex(name string) error { // shouldn't it also terminate the index' persistence loop goroutine?
	_, exists := mgr.indexes[name]
	if !exists {
		return fmt.Errorf("index '%s' does not exist", name)
	}
	delete(mgr.indexes, name)                                 // delete the index from memory
	return os.RemoveAll(path.Join(constants.EIGEN_DIR, name)) // delete the index from disk
}

func (mgr *IndexMgr) ListIndexes() ([]string, error) {
	_, indexNames, err := listPersistedIndexes()
	if err != nil {
		return nil, err
	}
	return indexNames, nil
}

func (mgr *IndexMgr) LoadIndexes(wg *sync.WaitGroup, e2eTestMode bool) error {
	defer wg.Done()
	if e2eTestMode { // skipping loading persisted indexes in E2E test mode, creating a test index instead
		fmt.Println("E2E test mode - skipping loading persisted indexes from disk.\nCreating E2E test index")
		return mgr.CreateIndex("e2e-test-index", 2, types.MetricL2)
	}

	// get persisted index data
	savedIndexesPaths, indexNames, err := listPersistedIndexes()
	if err != nil {
		return err
	}

	for i, indexPath := range savedIndexesPaths {
		if err := mgr.loadIndexFromDisk(indexPath, indexNames[i]); err != nil {
			return err
		}
	}
	return nil
}

func (mgr *IndexMgr) loadIndexFromDisk(indexPath string, indexName string) error {
	fmt.Printf("Loading index '%s' from disk...\n", indexName) // for testing

	idx := &index.Index{}
	indexDatafilePath := path.Join(constants.EIGEN_DIR, indexPath, constants.INDEX_DATAFILE)
	faissgoDatafilePath := path.Join(constants.EIGEN_DIR, indexPath, constants.FAISSGO_DATAFILE)

	// load the index struct (contains metadata, configuration values, etc...)
	serializedIndex, err := os.ReadFile(indexDatafilePath)
	if err != nil {
		return err
	}
	buf := bytes.NewBuffer(serializedIndex)
	decoder := gob.NewDecoder(buf)
	err = decoder.Decode(idx)
	if err != nil {
		return err
	}

	// load Faiss index (contains the embedding data)
	faissMetric, err := idx.Metric.ToFaissMetricType()
	if err != nil {
		return err
	}
	faissIdx, err := faissgo.IndexFactory(
		idx.Dimensions,
		constants.INDEX_TYPE_FAISS, // add more index types
		faissMetric,
	)
	if err != nil {
		return err
	}
	idx.SetFaissIndex(faissIdx)
	if err := mgr.loadFaissgoIndex(idx, faissgoDatafilePath); err != nil {
		return err
	}

	_, exists := mgr.indexes[indexName]
	if exists {
		return fmt.Errorf("index '%s' already exists in memory", indexName)
	}
	memIdxMgr.indexes[indexName] = idx
	return nil
}

func (mgr *IndexMgr) loadFaissgoIndex(index *index.Index, path string) error {
	return index.GetFaissIndex().LoadFromDisk(path)
}

// Persists an index to disk.
//
// The faissgo index is persisted separately as the index struct only stores
// a pointer to the faissgo index. This means that serializing the index struct
// will only serialize the pointer to the faissgo index, and not the faissgo index itself.
//
// Returns an error if one occured.
func (mgr *IndexMgr) persistIndex(index *index.Index) error {
	indexDatafilePath := path.Join(constants.EIGEN_DIR, index.Name, constants.INDEX_DATAFILE)
	faissgoDatafilePath := path.Join(constants.EIGEN_DIR, index.Name, constants.FAISSGO_DATAFILE)

	// persist the faissgo index
	if err := index.GetFaissIndex().WriteToDisk(faissgoDatafilePath); err != nil {
		return err
	}

	// serialize the index struct
	buf := new(bytes.Buffer)
	encoder := gob.NewEncoder(buf)
	err := encoder.Encode(index)
	if err != nil {
		return err
	}
	serializedData := buf.Bytes()
	return os.WriteFile(indexDatafilePath, serializedData, constants.DB_PERSIST_CHMOD)
}

// Starts a Goroutine that periodically persists all indexes to disk.
//
// The time interval at which data is persisted on disk is set
// in the "config" at persistence.timeInterval.
//
// Returns an error if one occured.
func (mgr *IndexMgr) startPersistenceLoop(wg *sync.WaitGroup) error {
	go func() {
		wg.Wait() // wait for all indexes to be loaded in memory before starting the persistence loop
		for {
			_, indexNames, err := listPersistedIndexes()
			if err != nil {
				fmt.Printf("Failed to list persisted indexes: %s\n", err.Error())
			}
			for _, indexName := range indexNames {
				idx, err := mgr.GetIndex(indexName)
				if err != nil {
					fmt.Printf("Index '%s' does not exist in memory, skipping persistence...\n", indexName)
					continue
				}
				if err := mgr.persistIndex(idx); err != nil {
					fmt.Printf("Failed to persist data to disk for index '%s': %s\n", indexName, err.Error())
				}
			}

			time.Sleep(cfg.GetConfig().GetPersistenceTimeInterval())
		}
	}()

	return nil
}

#include "lib/hnswlib/hnswlib.h"
#include "hnsw_wrapper.h"
#include <stdlib.h>

static thread_local std::string lastErrorMsg; // stores the last error message

/**
 * Returns the last error message if there exists one. If no error message exists, a null pointer will be returned.
 * 
 * @return last error message or a null pointer
 */
char* peekLastErrorMsg() {
    return lastErrorMsg.empty() ? nullptr : strdup(lastErrorMsg.c_str());
}

/**
 * Returns the last error message AND clears it, if there exists one. If no error message exists, a null pointer will be returned.
 * 
 * @return last error message or a null pointer
 */
char* getLastErrorMsg() {
    char *err = peekLastErrorMsg();
    lastErrorMsg.clear();
    return err;
}

/**
 * Instantiates and returns an HNSW index.
 *
 * @param dim:              dimension of the vector space
 * @param maxElements:      index's vector storage capacity
 * @param m:                `m` parameter in the HNSW algorithm
 * @param efConstruction:   `efConstruction` parameter in the HNSW algorithm
 * @param randSeed:         random seed
 * @param spaceType:        similarity metric to use in the index ("l" = L2, "i" = IP, "c" = cosine). (default: "l")
 * 
 * @return                  instance of a HNSW index
 */
HNSW initHNSW(int dim, unsigned long maxElements, int m, int efConstruction, int randSeed, char spaceType) {
    try {
        hnswlib::SpaceInterface<float> *vectorSpace;
        if (spaceType == 'i') { // inner product
            vectorSpace = new hnswlib::InnerProductSpace(dim);
        }
        else if (spaceType == 'c') { // cosine (cosine is the same as IP when all vectors are normalized)
            vectorSpace = new hnswlib::InnerProductSpace(dim);
        } else { // default: L2
            vectorSpace = new hnswlib::L2Space(dim);
        }
        return new hnswlib::HierarchicalNSW<float>(vectorSpace, maxElements, m, efConstruction, randSeed, true); // instantiate the hnsw index
    } catch (const std::runtime_error e) {
        lastErrorMsg = std::string(e.what());
        return nullptr;
    } catch(const std::exception e) {
        lastErrorMsg = std::string(e.what());
        return nullptr;
    }
}

/**
 * Generate an index using an index that has been saved on disk.
 * 
 * @param location:     the file path of the saved index
 * @param dim:          dimension of the vector space
 * @param spaceType:    similarity metric to use in the index ("l" = L2, "i" = IP, "c" = cosine). (default: "l")
 * @param maxElements:  index's vector storage capacity
 * 
 * @return              an index containing the data previously saved on disk
 */
HNSW loadHNSW(char *location, int dim, char spaceType, unsigned long maxElements) {
    try {
        hnswlib::SpaceInterface<float> *vectorSpace;
        if (spaceType == 'i') { // inner product
            vectorSpace = new hnswlib::InnerProductSpace(dim);
        }
        else if (spaceType == 'c') { // cosine (cosine is the same as IP when all vectors are normalized)
            vectorSpace = new hnswlib::InnerProductSpace(dim);
        } else { // default: L2
            vectorSpace = new hnswlib::L2Space(dim);
        }
        return new hnswlib::HierarchicalNSW<float>(vectorSpace, std::string(location), false, maxElements, true); // load the index from the specified location
    } catch (const std::runtime_error e) {
        lastErrorMsg = std::string(e.what());
        return nullptr;
    } catch(const std::exception e) {
        lastErrorMsg = std::string(e.what());
        return nullptr;
    }
}

/**
 * Saves an index as a file on the disk.
 * 
 * @param hnswIndex:    the HNSW index
 * @param location:     the location in which to save the index
 */
void saveHNSW(HNSW hnswIndex, char *location) {
    try {
        ((hnswlib::HierarchicalNSW<float>*) hnswIndex)->saveIndex(location);
    } catch (const std::exception e) {
        lastErrorMsg = std::string(e.what());
    }
}

/**
 * Frees an HNSW index from memory.
 *
 * @param hnswIndex: HNSW index to free
 */
void freeHNSW(HNSW hnswIndex) {
    delete (hnswlib::HierarchicalNSW<float>*) hnswIndex;
}

/**
 * Adds a vector to the HNSW index. 
 * NOTE: If a vector with the specified label already exists, IT WILL BE OVERWRITTEN.
 *
 * @param hnswIndex:    HNSW index to add the point to
 * @param vector:       the vector to add to the index
 * @param label:        the vector's label
 */
void insertVector(HNSW hnswIndex, float *vector, label_t label) {
    try {
        ((hnswlib::HierarchicalNSW<float>*) hnswIndex)->addPoint(vector, label, true);
    } catch(const std::runtime_error e) {
        lastErrorMsg = std::string(e.what());
    } catch (const std::exception e) {
        lastErrorMsg = std::string(e.what());
    }
}

/**
 * Returns a vector's components given its label.
 * 
 * @param hnswIndex:    the HNSW index
 * @param label:        the vector's label
 * 
 * @return              the vector's components
 */
float* getVector(HNSW hnswIndex, label_t label, int dim) {
    try {
        std::vector<float> vec = ((hnswlib::HierarchicalNSW<float>*) hnswIndex)->getDataByLabel<float>(label);
        float* vecPtr = (float*) malloc(sizeof(float) * dim);
        memcpy(vecPtr, vec.data(), sizeof(float) * dim);
        return vecPtr;
    } catch(const std::runtime_error e) {
        lastErrorMsg = std::string(e.what());
        return nullptr;
    } catch (const std::exception e) {
        lastErrorMsg = std::string(e.what());
        return nullptr;
    }
}

/**
 * Marks a vector with the specified label as deleted, which omits it from KNN search.
 * 
 * @param hnswIndex:    the HNSW index
 * @param label:        the vector's label
 */
void deleteVector(HNSW hnswIndex, label_t label) {
    try {
        ((hnswlib::HierarchicalNSW<float>*) hnswIndex)->markDelete(label);
    } catch(const std::runtime_error e) {
        lastErrorMsg = std::string(e.what());
    } catch (const std::exception e) {
        lastErrorMsg = std::string(e.what());
    }
}

/**
 * Performs similarity search on the HNSW index.
 * 
 * @param hnswIndex:    the HNSW index
 * @param vector:       the query vector
 * @param k:            the k value
 * @param labels:       a dynamic array which will receive the labels of the k-nearest neighbors
 * @param distances:    a dynamic array which will receive the distances of the k-nearest neighbors from the query vector
 * 
 * @return              the number of nearest neighbors found (num of nn <= k since it's possible for k > num of vectors in the index). If an error occured, it will return -1.
 */
int searchKNN(HNSW hnswIndex, float *vector, int k, label_t *labels, float *distances) {
    std::priority_queue<std::pair<float, hnswlib::labeltype>> searchResults;
    try {
        searchResults = ((hnswlib::HierarchicalNSW<float>*) hnswIndex)->searchKnn(vector, k); // use searchKnnCloserFirst to ensure the output is sorted

        int n = searchResults.size();
        std::pair<float, hnswlib::labeltype> pair;
        for (int i = n - 1; i >= 0; i--) {
            pair = searchResults.top();
            distances[i] = pair.first;
            labels[i] = pair.second;
            searchResults.pop();
        }
        return n;
    } catch (const std::runtime_error e) {
        lastErrorMsg = std::string(e.what());
        return -1;
    } catch (const std::exception e) {
        lastErrorMsg = std::string(e.what());
        return -1;
    }
}

/**
 * Set's the efConstruction parameter in the HNSW index.
 * 
 * @param hnswIndex:    the HNSW index
 * @param ef:           the new efConstruction parameter
 */
void setEf(HNSW hnswIndex, int ef) {
    ((hnswlib::HierarchicalNSW<float>*) hnswIndex)->ef_ = ef;
}
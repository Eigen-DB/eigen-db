package vectors

import "errors"

var vectorStorage []*vector

func InitializeVectorStorage() {
	vectorStorage = make([]*vector, 0)
}

func writeVectorToMemory(vector *vector) error {
	vectorStorage = append(vectorStorage, vector)
	return nil
}

func deleteVectorFromMemory(vectorId uint32) error {
	found := false
	for i, v := range vectorStorage {
		if v.id == vectorId {
			// removing the vector does not preserve the order
			vectorStorage[i] = vectorStorage[len(vectorStorage)-1]
			vectorStorage = vectorStorage[:len(vectorStorage)-1]
			found = true
			break
		}
	}

	if !found {
		return errors.New("vector not found in memory")
	}
	return nil
}

func exportMemoryToDisk() error { // TODO
	return nil
}

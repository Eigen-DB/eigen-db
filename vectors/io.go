package vectors

import (
	c "eigen_db/constants"
	"encoding/json"
	"os"
)

func writeVector(v *Vector) error {
	fi, err := os.Open(c.DATABASE_PATH)
	if err != nil {
		return err
	}
	defer fi.Close()

	vectorsByte := make([]byte, 0)
	_, err = fi.Read(vectorsByte)
	if err != nil {
		return err
	}

	vectors := make([]Vector, 0)
	err = json.Unmarshal(vectorsByte, &vectors)
	if err != nil {
		return err
	}

	vectors = append(vectors, *v)
	vectorsByte, err = json.Marshal(vectors)
	if err != nil {
		return err
	}

	_, err = fi.Write(vectorsByte)
	return err
}

func deleteVector(vectorId uint32) error { // TODO
	return nil
}

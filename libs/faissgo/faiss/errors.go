package faiss

/*
#cgo CXXFLAGS: -I${SRCDIR}/../lib/faiss/c_api
#cgo LDFLAGS: -lfaiss_c -lstdc++

#include <faiss/c_api/error_c.h>
*/
import "C"
import "errors"

func GetLastError() error {
	return errors.New(C.GoString(C.faiss_get_last_error()))
}

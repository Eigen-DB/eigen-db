package faissgo

/*
#include "lib/faiss/c_api/error_c.h"
*/
import "C"
import "errors"

func GetLastError() error {
	return errors.New(C.GoString(C.faiss_get_last_error()))
}

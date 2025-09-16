package constants

import (
	"io/fs"
)

const EIGEN_DIR string = "eigen"
const API_KEY_FILE_PATH string = EIGEN_DIR + "/api_key.txt"
const INDEX_DATAFILE_EXTENSION string = ".egn"
const INDEX_DATAFILE string = "index" + INDEX_DATAFILE_EXTENSION
const FAISSGO_DATAFILE string = "faissgo" + INDEX_DATAFILE_EXTENSION
const API_KEY_FILE_CHMOD fs.FileMode = 0600 // rw-------
const DB_PERSIST_CHMOD fs.FileMode = 0700   // rwx------ (x is needed to access files within the directory)
const CONFIG_CHMOD fs.FileMode = 0600       // rw-------
const CONFIG_PATH string = EIGEN_DIR + "/config.yml"
const TESTING_TMP_FILES_PATH string = "/tmp"
const ENV_VAR_API_KEY_NAME string = "EIGENDB_API_KEY"
const MIDDLEWARE_API_KEY_HEADER string = "X-Eigen-API-Key"
const INDEX_TYPE_FAISS string = "HNSW32,IDMap2"           // in the future, make this configurable per index
const VALID_INDEX_NAME_REGEX string = "^[a-z0-9-]{3,32}$" // only allow lowercase letters, numbers, and dashes, min length 3, max length 32

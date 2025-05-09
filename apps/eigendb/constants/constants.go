package constants

import (
	"io/fs"
)

const EIGEN_DIR string = "eigen"
const API_KEY_FILE_PATH string = EIGEN_DIR + "/api_key.txt"
const STORE_PERSIST_PATH string = EIGEN_DIR + "/eigendb-index.egn"
const INDEX_PERSIST_PATH string = EIGEN_DIR + "/faissgo-index.egn"
const API_KEY_FILE_CHMOD fs.FileMode = 0600 // rw-------
const DB_PERSIST_CHMOD fs.FileMode = 0600   // rw-------
const CONFIG_CHMOD fs.FileMode = 0600       // rw-------
const CONFIG_PATH string = EIGEN_DIR + "/config.yml"
const TESTING_TMP_FILES_PATH string = "/tmp"
const ENV_VAR_API_KEY_NAME string = "EIGENDB_API_KEY"
const MIDDLEWARE_API_KEY_HEADER string = "X-Eigen-API-Key"

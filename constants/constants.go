package constants

import "io/fs"

const EIGEN_DIR string = "./eigen"
const STORE_PERSIST_PATH string = EIGEN_DIR + "/vector_space.vec"
const INDEX_PERSIST_PATH string = EIGEN_DIR + "/hnsw_index.bin"
const DB_PERSIST_CHMOD fs.FileMode = 0600 // rw-------
const CONFIG_CHMOD fs.FileMode = 0600     // rw-------
const CONFIG_PATH string = EIGEN_DIR + "/config.yml"
const TESTING_TMP_FILES_PATH string = "/tmp"
const REDIS_API_KEY_NAME string = "apiKey"
const MIDDLEWARE_API_KEY_HEADER string = "X-Eigen-API-Key"

package constants

import "io/fs"

const EIGEN_DIR string = "./eigen"
const DB_PERSIST_FILENAME string = "vector_space.vec"
const DB_PERSIST_PATH string = EIGEN_DIR + "/" + DB_PERSIST_FILENAME
const DB_PERSIST_CHMOD fs.FileMode = 0777 // change
const CONFIG_PATH string = EIGEN_DIR + "/config.yml"
const TESTING_TMP_FILES_PATH string = "./test_utils/testing_tmp"

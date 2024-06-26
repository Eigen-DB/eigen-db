package vector

import (
	"eigen_db/types"
	"eigen_db/vector_io"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

const EXPECTED_RESPONSE_CODE int = 200
const EXPECTED_RESPONSE string = "Vector successfully inserted."

var newVectorInvokations int = 0
var insertInvokations int = 0

type MockVectorFactory struct{}

type MockVector struct{}

func (factory *MockVectorFactory) NewVector(components types.VectorComponents) vector_io.IVector {
	newVectorInvokations++
	return &MockVector{}
}

func (vector *MockVector) Insert() {
	insertInvokations++
}

func setupTestRouter() *gin.Engine {
	r := gin.Default()
	vectors := r.Group("/vector")
	vectors.PUT("/insert", Insert(&MockVectorFactory{}))
	return r
}

func TestInsertVector(t *testing.T) {
	router := setupTestRouter()
	w := httptest.NewRecorder()
	body := insertRequestBody{
		Components: types.VectorComponents{1.4, 2.3, 7.1},
	}
	bodyJSON, err := json.Marshal(body)
	if err != nil {
		t.Errorf("Error mashalling the test body as JSON: %s", err.Error())
	}
	req, err := http.NewRequest("PUT", "/vector/insert", strings.NewReader(string(bodyJSON)))
	if err != nil {
		t.Errorf("Error when creating (not sending) request to /vector/insert: %s", err.Error())
	}
	router.ServeHTTP(w, req)

	// checking for the proper response
	assert.Equal(t, EXPECTED_RESPONSE_CODE, w.Code)
	assert.Equal(t, EXPECTED_RESPONSE, w.Body.String())

	// checking to see if the functions within the handler were invoked the right amount of times
	assert.Equal(t, newVectorInvokations, 1)
	assert.Equal(t, insertInvokations, 1)
}

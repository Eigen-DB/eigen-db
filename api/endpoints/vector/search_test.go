package vector

import (
	"eigen_db/test_utils"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSearch(t *testing.T) {
	router := test_utils.SetupTestRouter()
	router.GET("/vector/search", Search(&test_utils.MockVectorSearcher{}))

	w := httptest.NewRecorder()
	body := searchRequestBody{
		QueryVectorId: uint32(1),
		K:             uint32(5),
	}
	bodyJSON, err := json.Marshal(body)
	if err != nil {
		t.Errorf("Error marshalling the test body as JSON: %s", err.Error())
	}
	req, err := http.NewRequest("GET", "/vector/search", strings.NewReader(string(bodyJSON)))
	if err != nil {
		t.Errorf("Error when creating (not sending) request to /vector/search: %s", err.Error())
	}
	router.ServeHTTP(w, req)

	// checking for the proper response
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "[1,2,3]", w.Body.String())

	// checking to see if the functions within the handler were invoked the right amount of times
	assert.Equal(t, test_utils.SimilaritySearchInvocations, 1)

	test_utils.Cleanup()
}

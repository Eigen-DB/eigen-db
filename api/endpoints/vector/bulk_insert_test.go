package vector

import (
	"eigen_db/test_utils"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type bulkInsertResponse struct {
	statusCode           int
	response             string
	insertInvocations    int
	newVectorInvocations int
}

func TestBulkInsert(t *testing.T) {
	router := test_utils.SetupTestRouter()
	router.PUT("/vector/bulk-insert", BulkInsert(&test_utils.MockVectorFactory{Dimensions: 3}))

	testcases := map[string]bulkInsertResponse{
		"{\"setOfComponents\": [[1.4, 2.3, 7.1], [3.4, 2.1, 3.4], [-5.2, 2.3, -8.7]]}": {
			statusCode:           200,
			response:             "3/3 vectors successfully inserted.",
			insertInvocations:    3,
			newVectorInvocations: 3,
		},
		"{\"setOfComponents\": []}": {
			statusCode:           200,
			response:             "0/0 vectors successfully inserted.",
			insertInvocations:    0,
			newVectorInvocations: 0,
		},
		"{\"setOfComponents\": [[1.4, 2.3], [3.4, 2.1, 3.4], [-5.2, 2.3, -8.7]]}": {
			statusCode:           200,
			response:             "2/3 vectors successfully inserted.",
			insertInvocations:    2,
			newVectorInvocations: 3,
		},
		"test": {
			statusCode:           400,
			response:             "",
			insertInvocations:    0,
			newVectorInvocations: 0,
		},
		"{}": {
			statusCode:           400,
			response:             "",
			insertInvocations:    0,
			newVectorInvocations: 0,
		},
		"": {
			statusCode:           400,
			response:             "",
			insertInvocations:    0,
			newVectorInvocations: 0,
		},
	}

	for input, output := range testcases {
		w := httptest.NewRecorder()
		req, err := http.NewRequest("PUT", "/vector/bulk-insert", strings.NewReader(input))
		if err != nil {
			t.Errorf("Error when creating (not sending) request: %s", err.Error())
		}
		router.ServeHTTP(w, req)

		// checking for the proper response
		assert.Equal(t, output.statusCode, w.Code, fmt.Sprintf("Body: %s", input))
		assert.Equal(t, output.response, w.Body.String(), fmt.Sprintf("Body: %s", input))

		// checking to see if the functions within the handler were invoked the right amount of times
		assert.Equal(t, output.insertInvocations, test_utils.InsertInvocations, fmt.Sprintf("Body: %s", input))
		assert.Equal(t, output.newVectorInvocations, test_utils.NewVectorInvocations, fmt.Sprintf("Body: %s", input))

		test_utils.Cleanup()
	}
}

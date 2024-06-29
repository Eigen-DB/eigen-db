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

type insertResponse struct {
	statusCode           int
	response             string
	insertInvocations    int
	newVectorInvocations int
}

func TestInsert(t *testing.T) {
	router := test_utils.SetupTestRouter()
	router.PUT("/vector/insert", Insert(&test_utils.MockVectorFactory{Dimensions: 3}))

	testcases := map[string]insertResponse{
		"{\"components\": [1.2, 3.5, -7.4]}": {
			statusCode:           200,
			response:             "Vector successfully inserted.",
			insertInvocations:    1,
			newVectorInvocations: 1,
		},
		"{\"components\": [1.2, 3.5]}": {
			statusCode:           400,
			response:             "vector provided had the wrong dimensionality",
			insertInvocations:    0,
			newVectorInvocations: 0,
		},
		"": {
			statusCode:           400,
			response:             "",
			insertInvocations:    0,
			newVectorInvocations: 0,
		},
		"{\"components\": []}": {
			statusCode:           400,
			response:             "vector provided had the wrong dimensionality",
			insertInvocations:    0,
			newVectorInvocations: 0,
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
	}

	for input, output := range testcases {
		w := httptest.NewRecorder()
		req, err := http.NewRequest("PUT", "/vector/insert", strings.NewReader(input))
		if err != nil {
			t.Errorf("Error when creating (not sending) the request: %s", err.Error())
		}
		router.ServeHTTP(w, req)

		// checking for the proper response
		assert.Equal(t, output.statusCode, w.Code, fmt.Sprintf("Body: %s", input))
		assert.Equal(t, output.response, w.Body.String(), fmt.Sprintf("Body: %s", input))

		// checking to see if the functions within the handler were invoked the right amount of times
		assert.Equal(t, output.insertInvocations, test_utils.InsertInvocations, fmt.Sprintf("Body: %s", input))
		assert.Equal(t, output.newVectorInvocations, test_utils.InsertInvocations, fmt.Sprintf("Body: %s", input))

		test_utils.Cleanup()
	}
}

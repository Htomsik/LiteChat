package Server

import (
	"Chat/internal/app/store/memoryStore"
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TODO realize canConnect tests with ws connection

// TestServer_SimpleEndpoint checking simple endpoints
func TestServer_SimpleEndpoint(t *testing.T) {

	// Arrange
	store := memoryStore.New()
	srv := newServer(store)

	cases := []struct {
		name         string
		body         interface{}
		expectedCode int
		httpMethod   string
		url          string
	}{
		{
			name:         "valid",
			body:         nil,
			expectedCode: http.StatusOK,
			httpMethod:   http.MethodGet,
			url:          "/api/isAlive",
		},
	}

	// Act
	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {

			recorder := httptest.NewRecorder()
			bodyBytes := &bytes.Buffer{}

			if testCase.body != nil {
				json.NewEncoder(bodyBytes).Encode(testCase.body)
			}

			request, _ := http.NewRequest(testCase.httpMethod, testCase.url, bodyBytes)

			srv.ServeHTTP(recorder, request)

			// Assert
			assert.Equal(t, testCase.expectedCode, recorder.Code)
		})
	}

}

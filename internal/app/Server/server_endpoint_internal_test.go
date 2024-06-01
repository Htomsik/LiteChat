package Server

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestServer_SimpleEndpoint checking simple endpoints
func TestServer_SimpleEndpoint(t *testing.T) {
	// Arrange
	srv := newServer()

	cases := []struct {
		name         string
		body         interface{}
		expectedCode int
		url          string
	}{
		{
			name:         "valid",
			body:         nil,
			expectedCode: http.StatusOK,
			url:          "/isAlive",
		},
	}

	// Act
	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()

			bodyBytes := &bytes.Buffer{}

			json.NewEncoder(bodyBytes).Encode(testCase.body)

			request, _ := http.NewRequest(http.MethodPost, testCase.url, bodyBytes)

			srv.ServeHTTP(recorder, request)

			// Assert
			assert.Equal(t, testCase.expectedCode, recorder.Code)
		})
	}

}

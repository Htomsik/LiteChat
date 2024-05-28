package apiServer

import (
	"Chat/internal/app/model"
	"Chat/internal/app/store/testStore"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

var (
	testCookieSecretKey = []byte("testSecretKey")
)

func TestServer_HandeUsersAdd(t *testing.T) {
	// Arrange
	srv := newServer(testStore.New(), sessions.NewCookieStore(testCookieSecretKey))

	testCases := []struct {
		name         string
		payload      interface{}
		expectedCode int
	}{
		{
			name: "valid",
			payload: map[string]string{
				"email":    "ex@ex.com",
				"password": "password",
			},
			expectedCode: http.StatusCreated,
		},
		{
			name:         "invalid payload",
			payload:      "nonePayload",
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "invalid payload params",
			payload: map[string]string{
				"email": "notEmail",
			},
			expectedCode: http.StatusUnprocessableEntity,
		},
	}

	// Act
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()

			bytesPayload := &bytes.Buffer{}
			json.NewEncoder(bytesPayload).Encode(testCase.payload)

			request, _ := http.NewRequest(http.MethodPost, userEndPoint, bytesPayload)

			srv.ServeHTTP(recorder, request)

			// Assert
			assert.Equal(t, testCase.expectedCode, recorder.Code)
		})
	}

}

func TestServer_HandeSessions(t *testing.T) {
	// Arrange
	user := model.TestUser(t)

	store := testStore.New()
	store.User().Add(user)

	srv := newServer(store, sessions.NewCookieStore(testCookieSecretKey))

	testCases := []struct {
		name         string
		payload      interface{}
		expectedCode int
	}{
		{
			name: "valid",
			payload: map[string]string{
				"email":    user.Email,
				"password": user.Password,
			},
			expectedCode: http.StatusOK,
		},
		{
			name:         "invalid payload",
			payload:      "nonePayload",
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "invalid payload params",
			payload: map[string]string{
				"email": "notEmail",
			},
			expectedCode: http.StatusUnauthorized,
		},
		{
			name: "invalid password",
			payload: map[string]string{
				"email":    user.Email,
				"password": "invalid",
			},
			expectedCode: http.StatusUnauthorized,
		},
	}

	// Act
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()

			bytesPayload := &bytes.Buffer{}
			json.NewEncoder(bytesPayload).Encode(testCase.payload)

			request, _ := http.NewRequest(http.MethodPost, userEndPoint+userAuthorize, bytesPayload)

			srv.ServeHTTP(recorder, request)

			// Assert
			assert.Equal(t, testCase.expectedCode, recorder.Code)
		})
	}

}

func TestServer_HandeAccountDeactivate(t *testing.T) {
	// Arrange
	user := model.TestUser(t)

	store := testStore.New()
	store.User().Add(user)

	userCookie := map[interface{}]interface{}{
		userIdSessionValue: user.ID,
	}

	srv := newServer(store, sessions.NewCookieStore(testCookieSecretKey))

	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodPut, accountActiveEndpoint+accountDeactivateEndpoint, nil)

	// Set encrypted auth cookie
	sc := securecookie.New(testCookieSecretKey, nil)
	cookie, _ := sc.Encode(sessionName, userCookie)
	request.Header.Set("Cookie", fmt.Sprintf("%s=%s", sessionName, cookie))

	// Act
	srv.ServeHTTP(recorder, request)

	// Assert
	assert.Equal(t, http.StatusOK, recorder.Code)
}

func TestServer_HandeAccountActivate(t *testing.T) {
	// Arrange
	user := model.TestUser(t)

	store := testStore.New()
	store.User().Add(user)
	store.User().Deactivate(user.ID)

	userCookie := map[interface{}]interface{}{
		userIdSessionValue: user.ID,
	}

	srv := newServer(store, sessions.NewCookieStore(testCookieSecretKey))

	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodPut, accountEndpoint+accountActivateEndpoint, nil)

	// Set encrypted auth cookie
	sc := securecookie.New(testCookieSecretKey, nil)
	cookie, _ := sc.Encode(sessionName, userCookie)
	request.Header.Set("Cookie", fmt.Sprintf("%s=%s", sessionName, cookie))

	// Act
	srv.ServeHTTP(recorder, request)

	// Assert
	assert.Equal(t, http.StatusOK, recorder.Code)
}

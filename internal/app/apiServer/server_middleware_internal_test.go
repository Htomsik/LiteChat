package apiServer

import (
	"Chat/internal/app/model"
	"Chat/internal/app/store/testStore"
	"fmt"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestServer_AuthenticateUser(t *testing.T) {
	// Arrange
	store := testStore.New()
	user := model.TestUser(t)

	store.User().Add(user)

	srv := newServer(store, sessions.NewCookieStore(testCookieSecretKey))

	// Create new cookie for response
	sc := securecookie.New(testCookieSecretKey, nil)
	dummyHandler := http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(http.StatusOK)
	})

	testCases := []struct {
		name         string
		cookieValue  map[interface{}]interface{}
		expectedCode int
	}{
		{
			name: "authenticated",
			cookieValue: map[interface{}]interface{}{
				userIdSessionValue: user.ID,
			},
			expectedCode: http.StatusOK,
		},
		{
			name:         "not authenticated",
			cookieValue:  nil,
			expectedCode: http.StatusUnauthorized,
		},
	}

	// Act
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()

			request, _ := http.NewRequest(http.MethodGet, homeEndpoint, nil)

			// Generate cookie session key for response
			cookie, _ := sc.Encode(sessionName, testCase.cookieValue)
			request.Header.Set("Cookie", fmt.Sprintf("%s=%s", sessionName, cookie))

			srv.authenticateUserMiddleWare(dummyHandler).ServeHTTP(recorder, request)

			// Assert
			assert.Equal(t, testCase.expectedCode, recorder.Code)
		})
	}
}

func TestServer_ActiveUser(t *testing.T) {
	// Arrange
	store := testStore.New()
	ActiveUser := model.TestUser(t)
	NotActiveUser := model.TestUser(t)

	NotActiveUser.Email = "NotActiveUser@ex.com"

	store.User().Add(ActiveUser)
	store.User().Add(NotActiveUser)

	store.User().Deactivate(NotActiveUser.ID)

	srv := newServer(store, sessions.NewCookieStore(testCookieSecretKey))

	// Create new cookie for response
	sc := securecookie.New(testCookieSecretKey, nil)
	dummyHandler := http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(http.StatusOK)
	})

	testCases := []struct {
		name         string
		cookieValue  map[interface{}]interface{}
		expectedCode int
	}{
		{
			name: "Active",
			cookieValue: map[interface{}]interface{}{
				userIdSessionValue: ActiveUser.ID,
			},
			expectedCode: http.StatusOK,
		},
		{
			name: "not Active",
			cookieValue: map[interface{}]interface{}{
				userIdSessionValue: NotActiveUser.ID,
			},
			expectedCode: http.StatusUnauthorized,
		},
		{
			name:         "nil",
			cookieValue:  nil,
			expectedCode: http.StatusUnauthorized,
		},
		{
			name: "notExist",
			cookieValue: map[interface{}]interface{}{
				userIdSessionValue: -1,
			},
			expectedCode: http.StatusUnauthorized,
		},
	}

	// Act
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()

			request, _ := http.NewRequest(http.MethodGet, accountEndpoint+accountActiveEndpoint+accountWhoAmIEndpoint, nil)

			// Generate cookie session key for response
			cookie, _ := sc.Encode(sessionName, testCase.cookieValue)
			request.Header.Set("Cookie", fmt.Sprintf("%s=%s", sessionName, cookie))

			srv.authenticateUserMiddleWare(srv.activeUserMiddleWare(dummyHandler)).ServeHTTP(recorder, request)

			// Assert
			assert.Equal(t, testCase.expectedCode, recorder.Code)
		})
	}
}

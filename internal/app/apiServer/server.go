package apiServer

import (
	"Chat/internal/app/store"
	"encoding/json"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/sirupsen/logrus"
	"net/http"
)

const (
	homeEndpoint = "/"
)

// server ...
type server struct {
	router       *mux.Router
	logger       *logrus.Logger
	store        store.Store
	sessionStore sessions.Store
}

// newServer ...
func newServer(store store.Store, sessionStore sessions.Store) *server {
	srv := &server{
		router:       mux.NewRouter(),
		logger:       logrus.New(),
		store:        store,
		sessionStore: sessionStore,
	}

	srv.configureRouter()

	return srv
}

func (srv *server) configureRouter() {
	// Add access with different domains
	srv.router.Use(handlers.CORS(handlers.AllowedOrigins([]string{"*"})))

	// Add requestID for all endpoints
	srv.router.Use(srv.requestIDMiddleWare)

	// Add logger middleware for all endpoints
	srv.router.Use(srv.logRequestMiddleWare)

	srv.configureOtherEndpoints()
	srv.configureAccountEndpoint()
	srv.configureAccountActiveEndpoints()
	srv.configureUserEndpoint()
}

func (srv *server) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	srv.router.ServeHTTP(writer, request)
}

// error called if there are any errors in the request
func (srv *server) error(writer http.ResponseWriter, request *http.Request, code int, err error) {
	srv.respond(writer, request, code, map[string]string{"error": err.Error()})
}

// respond on request
func (srv *server) respond(writer http.ResponseWriter, request *http.Request, code int, data interface{}) {
	writer.WriteHeader(code)

	if data != nil {
		json.NewEncoder(writer).Encode(data)
	}
}

package apiServer

import (
	_ "Chat/docs"
	_ "github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
	"net/http"
)

// Account endpoints
const (
	accountEndpoint         = "/account"
	accountActivateEndpoint = "/activate"
)

// users endpoints
const (
	userEndPoint  = "/user"
	userAuthorize = "/authorize"
)

// Account/Active endpoints
const (
	accountActiveEndpoint     = accountEndpoint + "/active"
	accountWhoAmIEndpoint     = "/who"
	accountDeactivateEndpoint = "/deactivate"
)

// configureOtherEndpoints public endpoints
func (srv *server) configureOtherEndpoints() {
	srv.router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)
}

// configureAccountEndpoint Account endpoints with authentication middleware
func (srv *server) configureAccountEndpoint() {
	account := srv.router.PathPrefix(accountEndpoint).Subrouter()
	account.Use(srv.authenticateUserMiddleWare)

	account.HandleFunc(accountActivateEndpoint, srv.handleAccountActivate()).Methods(http.MethodPut)
}

// configureUserEndpoint user endpoints for create/authorize
func (srv *server) configureUserEndpoint() {
	user := srv.router.PathPrefix(userEndPoint).Subrouter()

	user.HandleFunc("", srv.handleUserAdd()).Methods(http.MethodPost)
	user.HandleFunc(userAuthorize, srv.handleUserSessionAdd()).Methods(http.MethodPost)
}

// configureAccountActiveEndpoints Account endpoints with authentication + active middleware
func (srv *server) configureAccountActiveEndpoints() {
	accountActive := srv.router.PathPrefix(accountActiveEndpoint).Subrouter()
	accountActive.Use(srv.authenticateUserMiddleWare)
	accountActive.Use(srv.activeUserMiddleWare)

	accountActive.HandleFunc(accountWhoAmIEndpoint, srv.handleWho()).Methods(http.MethodGet)
	accountActive.HandleFunc(accountDeactivateEndpoint, srv.handleAccountDeactivate()).Methods(http.MethodPut)
}

package Server

import (
	_ "Chat/docs"
	_ "github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
)

// configureEndpoints endpoint configurations
func (srv *server) configureEndpoints() {

	srv.configureApiEndpoint()
	srv.configurePageEndpoint()
}

// configureApiEndpoint internal functional endpoints
func (srv *server) configureApiEndpoint() {
	apiRouter := srv.router.PathPrefix("/api").Subrouter()

	apiRouter.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)
	apiRouter.HandleFunc("/isAlive", srv.handleIsAlive())

	apiRouter.HandleFunc("/chat/{id}", srv.handleChat())
}

// configurePageEndpoint html pages
func (srv *server) configurePageEndpoint() {

	srv.router.HandleFunc("/chat/{id}", srv.handleHomePage())
}

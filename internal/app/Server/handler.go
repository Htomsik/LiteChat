package Server

import (
	_ "Chat/docs"
	_ "github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
)

// configureEndpoints endoint configurations
func (srv *server) configureEndpoints() {

	srv.configureFunctionalEndpoint()
}

// configureFunctionalEndpoint internal functional endpoints
func (srv *server) configureFunctionalEndpoint() {
	srv.router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)
	srv.router.HandleFunc("/isAlive", srv.handleIsAlive())
}

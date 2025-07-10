package Server

import (
	_ "Chat/docs"
	"Chat/internal/app/model"

	_ "github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
)

// configureEndpoints endpoint configurations
func (srv *server) configureEndpoints() {

	srv.configureApiEndpoint()
	srv.configureChatRouter()
	srv.configurePageEndpoint()
}

// configureApiEndpoint internal functional endpoints
func (srv *server) configureApiEndpoint() {
	apiRouter := srv.router.PathPrefix("/api").Subrouter()

	apiRouter.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)
	apiRouter.HandleFunc("/isAlive", srv.handleIsAlive())

}

func (srv *server) configureChatRouter() {
	chatRouter := srv.router.PathPrefix("/api/chat").Subrouter()

	chatRouter.Use(srv.chatUserMiddleWare)

	chatRouter.HandleFunc("/{id}", srv.handleChat())
	chatRouter.HandleFunc("/canConnect/{id}", srv.handleCanConnect())
}

// configurePageEndpoint html pages
func (srv *server) configurePageEndpoint() {

	websiteHandler := model.SpaHandler{
		StaticPath: "website/dist",
		IndexPath:  "index.html",
	}

	srv.router.PathPrefix("/").Handler(websiteHandler)
}

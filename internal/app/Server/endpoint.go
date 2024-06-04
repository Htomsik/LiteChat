package Server

import (
	"Chat/internal/app/model"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"net/http"
)

// TODO придумать как по нормальному это хранить
var hubs = make(map[string]*model.Hub)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// handleIsAlive Checking server is alive
// @Summary      checking server is alive
// @Success      200
// @Router       /api/isAlive [Get]
func (srv *server) handleIsAlive() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		srv.respond(w, r, http.StatusOK, nil)
	}
}

// handleCanConnect Checking can connect to server
// @Summary      Checking can connect to server
// @Success      200
// @Router       /api/canConnect [get]
func (srv *server) handleCanConnect() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Check chat id
		pathVars := mux.Vars(r)
		hubId, ok := pathVars["id"]

		if !ok {
			srv.respond(w, r, http.StatusBadRequest, nil)
			return
		}

		// Check userQuery
		userQuery := r.URL.Query().Get(model.QueryValueUser)
		if userQuery == "" {
			srv.error(w, r, http.StatusBadRequest, errors.New(fmt.Sprintf(model.QueryVariableNotFound, model.QueryValueUser)))
			return
		}

		// if chat doesn't exists
		hub, ok := hubs[hubId]
		if !ok {
			srv.respond(w, r, http.StatusOK, nil)
			return
		}

		// Check is user with same name is connected
		user := r.Context().Value(contextUser).(*model.ChatUser)
		if hub.FindClient(user.Name) != nil {
			srv.respond(w, r, http.StatusUnprocessableEntity, fmt.Sprintf(model.AlreadyExists, user.Name))
			return
		}

		srv.respond(w, r, http.StatusOK, nil)
	}
}

// handleChat websocket chat
// @Summary      Connecting to websocket chat
// @Success      200
// @Param        id   path      string  true  "Chat id"
// @Router       /api/chat/{id} [Get]
func (srv *server) handleChat() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		// Check chat id
		pathVars := mux.Vars(r)
		hubId, ok := pathVars["id"]

		if !ok {
			srv.respond(w, r, http.StatusBadRequest, nil)
			return
		}

		// Check is chat exists by id
		hub, ok := hubs[hubId]

		// is no exists create new
		if !ok {
			hub = model.HewHub(srv.logger)
			hubs[hubId] = hub
			go hub.Run()
		}

		// Check is user with same name is connected
		user := r.Context().Value(contextUser).(*model.ChatUser)
		if hub.FindClient(user.Name) != nil {
			srv.respond(w, r, http.StatusUnprocessableEntity, fmt.Sprintf(model.AlreadyExists, user.Name))
			return
		}

		// Create new websocket connection
		connection, err := upgrader.Upgrade(w, r, nil)

		if err != nil {
			srv.logger.Infof("Can't create websocket connection: %v", err)
			return
		}

		client := model.NewClient(hub, connection, user)
		client.RegisterToHub()

		go client.WriteToHub()
		go client.ReadFromHub()
	}
}

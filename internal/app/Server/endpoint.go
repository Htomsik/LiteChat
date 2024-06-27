package Server

import (
	"Chat/internal/app/model/chat"
	"Chat/internal/app/model/constants"
	client "Chat/internal/app/model/websocket"
	"errors"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"net/http"
)

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
// @Router       /api/canConnect/{id} [get]
func (srv *server) handleCanConnect() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Check chat id
		pathVars := mux.Vars(r)
		hubId, ok := pathVars["id"]

		if !ok {
			srv.respond(w, r, http.StatusBadRequest, nil)
			return
		}

		// if chat doesn't exists
		hub, err := srv.store.Hub().Find(hubId)

		// if error is not record not found
		if err != nil && !errors.Is(err, constants.ErrorRecordNotFound) {
			srv.respond(w, r, http.StatusInternalServerError, err)
			return
		}

		if hub == nil {
			srv.respond(w, r, http.StatusOK, nil)
			return
		}

		// Check is all right with user
		user := r.Context().Value(contextUser).(*chat.User)
		if user == nil {
			srv.respond(w, r, http.StatusInternalServerError, "")
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
		hub, err := srv.store.Hub().Find(hubId)

		// if error is not record not found
		if err != nil && !errors.Is(err, constants.ErrorRecordNotFound) {
			srv.respond(w, r, http.StatusInternalServerError, err)
			return
		}

		// is no exists create new
		newHub := false
		if hub == nil {
			hub, err = srv.store.Hub().Create(hubId)
			if err != nil {
				srv.respond(w, r, http.StatusInternalServerError, nil)
			}
			newHub = true
			go hub.Run()
		}

		user := r.Context().Value(contextUser).(*chat.User)
		if user == nil {
			srv.respond(w, r, http.StatusInternalServerError, err)
			return
		}

		// Create new websocket connection
		connection, err := upgrader.Upgrade(w, r, nil)

		if err != nil {
			srv.logger.Infof("Can't create websocket connection: %v", err)
			return
		}

		client := client.NewClient(hub.Commands, connection, user)
		client.RegisterToHub()

		// if hub is new first person is admin
		if newHub {
			hub.SetAdmin(client)
		}

		go client.WriteToHub()
		go client.ReadFromHub()
	}
}

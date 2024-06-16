package Server

import (
	"Chat/internal/app/model"
	"errors"
	"fmt"
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
		if err != nil && !errors.Is(err, model.ErrorRecordNotFound) {
			srv.respond(w, r, http.StatusInternalServerError, err)
			return
		}

		if hub == nil {
			srv.respond(w, r, http.StatusOK, nil)
			return
		}

		// Check is all right with user
		user := r.Context().Value(contextUser).(*model.ChatUser)
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
		if err != nil && !errors.Is(err, model.ErrorRecordNotFound) {
			srv.respond(w, r, http.StatusInternalServerError, err)
			return
		}

		// is no exists create new
		if hub == nil {
			hub = model.HewHub(hubId, srv.logger)

			err = srv.store.Hub().Add(hub)
			if err != nil {
				srv.respond(w, r, http.StatusInternalServerError, nil)
			}

			go hub.Run()
		}

		// Check is user with same originalName is connected
		// if yes change name +1
		user := r.Context().Value(contextUser).(*model.ChatUser)
		if count := hub.CountUsersByOriginalName(user.Name); count > 0 {
			user.Name = fmt.Sprintf("%v[%v]", user.Name, count)
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

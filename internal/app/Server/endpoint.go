package Server

import (
	"Chat/internal/app/model/Hub"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"net/http"
)

// TODO придумать как по нормальному это хранить
var hubs = make(map[string]*Hub.Hub)

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

// handleHomePage домашняя страницы
func (srv *server) handleHomePage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./web/static/home.html")
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
		}

		// Проверяем создан ли чат по такой ссылке
		hub, ok := hubs[hubId]

		// Если нет создаем новый
		if !ok {
			hub = Hub.HewHub(srv.logger)
			hubs[hubId] = hub
			go hub.Run()
		}

		// Создаем новое соединение вебсокет
		connection, err := upgrader.Upgrade(w, r, nil)

		if err != nil {
			srv.logger.Infof("Can't create websocket connection: %v", err)
			return
		}

		// Создаем нового клиента для этого соединения
		client := Hub.NewClient(hub, connection)

		client.RegisterToHub()

		go client.WriteToHub()
		go client.ReadFromHub()
	}
}

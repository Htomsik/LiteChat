package Server

import (
	"net/http"
)

// handleIsAlive Checking server is alive
// @Summary      checking server is alive
// @Success      200
// @Router       /isAlive [Get]
func (srv *server) handleIsAlive() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		srv.respond(w, r, http.StatusOK, nil)
	}
}

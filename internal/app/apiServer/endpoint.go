package apiServer

import (
	"Chat/internal/app/model"
	"encoding/json"
	"net/http"
)

const (
	sessionName        = "Chat"
	userIdSessionValue = "userId"
)

// handleUserAdd Create User
// @Summary      Add account/User
// @Description  Create new account
// @Tags         User
// @Accept       json
// @Produce      json
// @Param		 User 	body 		model.UserShort 	true 	"user information"
// @Success      201  	{object} 	model.User
// @Failure      422
// @Failure      400
// @Router       /user [post]
func (srv *server) handleUserAdd() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Decode request
		req := &model.UserShort{}

		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			srv.error(w, r, http.StatusBadRequest, err)
			return
		}

		// Add new user
		user := &model.User{
			Email:    req.Email,
			Password: req.Password,
		}

		if err := srv.store.User().Add(user); err != nil {
			srv.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		user.ClearPrivate()
		srv.respond(w, r, http.StatusCreated, user)
	}
}

// handleUserSessionAdd Authorize into User
// @Summary      Authorize into account
// @Description  Authorize into account by session cookie
// @Tags         User
// @Accept       json
// @Produce      json
// @Param		 User 	body	model.UserShort 	true 	"user information"
// @Failure      401
// @Failure      400
// @Failure      500
// @Success      200
// @Router       /user/authorize [post]
func (srv *server) handleUserSessionAdd() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Decode request
		req := &model.UserShort{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			srv.error(w, r, http.StatusBadRequest, err)
			return
		}

		// Find user
		user, err := srv.store.User().FindByEmail(req.Email)
		if err != nil || !user.ComparePassword(req.Password) {
			srv.error(w, r, http.StatusUnauthorized, model.EmailOrPasswordIncorrect)
			return
		}

		// Create user session
		session, err := srv.sessionStore.Get(r, sessionName)
		if err != nil {
			srv.error(w, r, http.StatusInternalServerError, err)
		}

		session.Values[userIdSessionValue] = user.ID

		// Add user session into sessions store
		if err = srv.sessionStore.Save(r, w, session); err != nil {
			srv.error(w, r, http.StatusInternalServerError, err)
		}

		srv.respond(w, r, http.StatusOK, nil)
	}
}

// handleWho info about current authorised user
// @Summary      Account info
// @Description  info about current user
// @Tags         Account
// @Accept       json
// @Produce      json
// @Success      200	{object}	model.User
// @Failure      401
// @Router       /account/active/who [get]
func (srv *server) handleWho() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		srv.respond(writer, request, http.StatusOK, request.Context().Value(ctxKeyUser).(*model.User))
	}
}

// handleAccountDeactivate deactivate current active account
// @Summary      Deactivate account
// @Description  Only deactivate, not delete
// @Tags         Account
// @Accept       json
// @Produce      json
// @Success      200
// @Failure      422
// @Failure      401
// @Router       /account/active/deactivate [put]
func (srv *server) handleAccountDeactivate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value(ctxKeyUser).(*model.User)

		if err := srv.store.User().Deactivate(user.ID); err != nil {
			srv.error(w, r, http.StatusUnprocessableEntity, err)
		}

		srv.respond(w, r, http.StatusOK, nil)
	}
}

// handleAccountActivate activate nonactive account
// @Summary      Activate account
// @Description  Activate only deactivated accounts
// @Tags         Account
// @Accept       json
// @Produce      json
// @Success      200
// @Failure      422
// @Router       /account/activate [put]
func (srv *server) handleAccountActivate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value(ctxKeyUser).(*model.User)

		if err := srv.store.User().Activate(user.ID); err != nil {
			srv.error(w, r, http.StatusUnprocessableEntity, err)
		}

		srv.respond(w, r, http.StatusOK, nil)
	}
}

package apiServer

import (
	"Chat/internal/app/model"
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

const (
	ctxKeyUser ctxKey = iota
	ctxRequestId
)

const (
	requestIdHeader = "Request-ID"
)

type ctxKey int8

func (srv *server) authenticateUserMiddleWare(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {

		// Get user session from request
		session, err := srv.sessionStore.Get(request, sessionName)
		if err != nil {
			srv.error(writer, request, http.StatusInternalServerError, err)
			return
		}

		// Get user id from session
		userId, exist := session.Values[userIdSessionValue]
		if !exist {
			srv.error(writer, request, http.StatusUnauthorized, model.NotAuthenticated)
			return
		}

		// Check user id exist
		user, err := srv.store.User().Find(userId.(int))
		if err != nil {
			srv.error(writer, request, http.StatusUnauthorized, model.NotAuthenticated)
			return
		}

		// Throw user context next
		next.ServeHTTP(writer, request.WithContext(context.WithValue(request.Context(), ctxKeyUser, user)))
	})
}

func (srv *server) activeUserMiddleWare(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {

		user := request.Context().Value(ctxKeyUser)

		if user == nil {
			srv.error(writer, request, http.StatusInternalServerError, model.ContextNotFound)
			return
		}

		if !user.(*model.User).Active {
			srv.error(writer, request, http.StatusUnauthorized, model.NotActive)
			return
		}

		next.ServeHTTP(writer, request)
	})
}

func (srv *server) requestIDMiddleWare(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		// Generate new guid
		guid := uuid.New().String()

		// Set guid to header
		writer.Header().Set(requestIdHeader, guid)

		// Throw request id next
		next.ServeHTTP(writer, request.WithContext(context.WithValue(request.Context(), ctxRequestId, guid)))
	})
}

func (srv *server) logRequestMiddleWare(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		// Create local logger rules
		logger := srv.logger.WithFields(logrus.Fields{
			"remove_addr": request.RemoteAddr, // request address
			"request_id":  request.Context().Value(ctxRequestId),
		})

		requestInfo := fmt.Sprintf("%s %s", request.Method, request.RequestURI)

		logger.Infof("Attempt %v", requestInfo)
		startTime := time.Now()

		customWriter := &responseWriter{writer, http.StatusOK}

		next.ServeHTTP(customWriter, request)

		loggerText := fmt.Sprintf(
			"Completed %v with %d %s in %v",
			requestInfo,
			customWriter.code,
			http.StatusText(customWriter.code),
			time.Now().Sub(startTime))

		if customWriter.code == http.StatusInternalServerError {
			logger.Error(loggerText)
		} else {
			logger.Infof(loggerText)
		}
	})
}

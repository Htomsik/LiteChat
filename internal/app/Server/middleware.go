package Server

import (
	"Chat/internal/app/model"
	"Chat/internal/app/model/chat"
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"net/http"
	"regexp"
	"strings"
	"time"
)

const (
	contextUser ctxKey = iota
	contextRequestId
)

type ctxKey int8

// chatUserMiddleWare throw userInfo from session to context
func (srv *server) chatUserMiddleWare(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Check is userNameMatch header exists
		userQuery := r.URL.Query().Get(model.QueryValueUser)
		if userQuery == "" {
			srv.error(w, r, http.StatusBadRequest, errors.New(fmt.Sprintf(model.QueryVariableNotFound, model.QueryValueUser)))
			return
		}

		// Check is userName valid
		userName := strings.TrimSpace(userQuery)
		if len(userName) < 2 || len(userName) > 20 {
			srv.error(w, r, http.StatusUnprocessableEntity, errors.New(fmt.Sprintf("%v user must be in the range of 2 to 20 characters", model.QueryValueUser)))
			return
		}

		var re = regexp.MustCompile(`[a-zA-Z0-9]`)
		userNameMatch := re.ReplaceAllString(userName, "")

		if userNameMatch != "" {
			srv.error(w, r, http.StatusUnprocessableEntity, errors.New(fmt.Sprintf("%v must be only from numbers and latin symbols", model.QueryValueUser)))
			return
		}

		chatUser := chat.NewChatUser(userName)

		// Throw userNameMatch context next
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), contextUser, chatUser)))
	})
}

func (srv *server) requestIDMiddleWare(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		// Generate new guid
		guid := uuid.New().String()

		// Set guid to header
		writer.Header().Set(model.RequestIdHeader, guid)

		// Throw request id next
		next.ServeHTTP(writer, request.WithContext(context.WithValue(request.Context(), contextRequestId, guid)))
	})
}

func (srv *server) logRequestMiddleWare(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		// Create local logger rules
		logger := srv.logger.WithFields(logrus.Fields{
			"remove_addr": request.RemoteAddr, // request address
			"request_id":  request.Context().Value(contextRequestId),
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

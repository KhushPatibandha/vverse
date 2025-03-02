package middleware

import (
	"errors"
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/KhushPatibandha/vverse/api"
)

var UnAuthError = errors.New("Invalid token.")

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// > All API calls must be authenticated (assume static API tokens)
		// In case of actual implementation, connect to the database to verify the token, also use JWT
		token := r.Header.Get("Authorization")
		if token == "" {
			log.Error(UnAuthError)
			api.RequestErrorHandler(w, UnAuthError)
			return
		}

		// simulate fetching token from database
		// time.Sleep(time.Second * 2)

		if token != "someCrazySecureToken" {
			log.Error(UnAuthError)
			api.RequestErrorHandler(w, UnAuthError)
			return
		}
		next.ServeHTTP(w, r)
	})
}

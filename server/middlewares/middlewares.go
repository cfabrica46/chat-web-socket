package middlewares

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/cfabrica46/chat-web-socket/server/database"
	"github.com/cfabrica46/chat-web-socket/server/token"
	"github.com/gorilla/mux"
)

type ErrMessage struct {
	Message string
}

type ContextKey string

var ContextUserKey ContextKey = "data-user"
var ContextMessageKey ContextKey = "data-message"

func GetUser() mux.MiddlewareFunc {

	return func(next http.Handler) http.Handler {

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			var errMessage ErrMessage

			tokenString := r.Header.Get("Authorization-header")

			check := database.CheckIfTokenIsInBlackList(tokenString, database.DB)

			if !check {
				errMessage.Message = "El token no es válido"
				json.NewEncoder(w).Encode(errMessage)
				return
			}

			user, err := token.ExtractUserFromClaims(tokenString)

			if err != nil {
				errMessage.Message = err.Error()
				json.NewEncoder(w).Encode(errMessage)
				return
			}

			user.Token = tokenString

			deadline, err := time.Parse(time.ANSIC, user.Deadline)

			if err != nil {
				errMessage.Message = http.StatusText(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(errMessage)
				return
			}

			checkTime := time.Now().Local().After(deadline)

			if !checkTime {
				errMessage.Message = "El token no es válido"
				json.NewEncoder(w).Encode(errMessage)
				return
			}

			ctx := context.WithValue(r.Context(), ContextUserKey, user)

			next.ServeHTTP(w, r.WithContext(ctx))

		})
	}

}

func LoginPassword() mux.MiddlewareFunc {

	return func(next http.Handler) http.Handler {

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			var errMessage ErrMessage
			var user database.User

			err := json.NewDecoder(r.Body).Decode(&user)

			if err != nil {
				errMessage.Message = http.StatusText(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(errMessage)
				return
			}

			ctx := context.WithValue(r.Context(), ContextUserKey, user)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}

}

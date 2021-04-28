package middlewares

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/cfabrica46/chat-web-socket/server/database"
	"github.com/cfabrica46/chat-web-socket/server/token"
)

type ErrMessage struct {
	Message string
}

type ContextKey string

var ContextUserKey ContextKey

func GetUser(next http.Handler, db *sql.DB) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var errMessage ErrMessage

		tokenString := r.Header.Get("Authorization-header")

		check := database.CheckIfTokenIsInBlackList(tokenString, db)

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

		user.Token = tokenString

		ctx := context.WithValue(r.Context(), ContextUserKey, user)

		next.ServeHTTP(w, r.WithContext(ctx))

	})

}

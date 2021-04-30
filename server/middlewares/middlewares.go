package middlewares

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"strings"
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

func GetUser(db *sql.DB) mux.MiddlewareFunc {

	return func(next http.Handler) http.Handler {

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			var errMessage ErrMessage

			URL := r.URL.String()

			endPoitns := strings.Split(URL, "/")

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

			switch endPoitns[1] {
			case "logout":

				stmt, err := db.Prepare("INSERT INTO black_list(token) VALUES (?)")

				if err != nil {
					errMessage.Message = http.StatusText(http.StatusInternalServerError)
					json.NewEncoder(w).Encode(errMessage)
					return
				}

				_, err = stmt.Exec(user.Token)

				if err != nil {
					errMessage.Message = http.StatusText(http.StatusInternalServerError)
					json.NewEncoder(w).Encode(errMessage)
					return
				}

				next.ServeHTTP(w, r)

			case "user":

				ctx := context.WithValue(r.Context(), ContextUserKey, user)

				next.ServeHTTP(w, r.WithContext(ctx))

			case "chat":

				ctx := context.WithValue(r.Context(), ContextUserKey, user)

				next.ServeHTTP(w, r.WithContext(ctx))

			}

		})
	}

}

func LoginPassword(db *sql.DB) mux.MiddlewareFunc {

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

			switch r.URL.String() {
			case "/login":

				row := db.QueryRow("SELECT users.id,users.role FROM users WHERE users.username=? AND users.password=?", user.Username, user.Password)

				err := row.Scan(&user.ID, &user.Role)

				if err != nil {
					errMessage.Message = http.StatusText(http.StatusBadRequest)
					json.NewEncoder(w).Encode(errMessage)
					return
				}

			case "/register":

				err := database.AddUser(&user, db)

				if err != nil {
					errMessage.Message = err.Error()
					json.NewEncoder(w).Encode(errMessage)
					return
				}

			}
			ctx := context.WithValue(r.Context(), ContextUserKey, user)
			next.ServeHTTP(w, r.WithContext(ctx))

		})
	}

}

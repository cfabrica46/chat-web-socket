package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/cfabrica46/chat-web-socket/server/database"
	"github.com/cfabrica46/chat-web-socket/server/token"
	"github.com/gorilla/mux"
)

var db database.DB

func main() {

	d, err := database.Open()

	if err != nil {
		log.Fatal(err)
	}

	db = database.DB{
		D: d,
	}

	go database.CleanBlackList(db.D)

	r := mux.NewRouter()

	r.HandleFunc("/login", login)
	r.HandleFunc("/register", register)

	s := r.PathPrefix("/").Subrouter()
	s.Use(middlewareGetUser)

	s.HandleFunc("/user", user)
	s.HandleFunc("/logout", logout)

	fmt.Println("Listening on localhost:8080")

	http.ListenAndServe(":8080", r)
}

func middlewareGetUser(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var errMessage ErrMessage

		tokenString := r.Header.Get("Authorization-header")

		check := database.CheckIfTokenIsInBlackList(tokenString, db.D)

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

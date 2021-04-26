package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type User struct {
	ID                                        int
	Username, Password, Deadline, Role, Token string
}

type ErrMessage struct {
	Message string
}

type ContextKey string

var ContextUserKey ContextKey

var db *sql.DB

func main() {

	var err error

	r := mux.NewRouter()

	db, err = open()

	if err != nil {
		log.Fatal(err)
	}

	go cleanBlackList()

	r.HandleFunc("/user", user).Subrouter().Use(middlewareGetUser)
	r.HandleFunc("/login", login)
	r.HandleFunc("/register", register)
	r.HandleFunc("/loout", logout).Subrouter().Use(middlewareGetUser)

	fmt.Println("Listening on localhost:8080")

	http.ListenAndServe(":8080", r)
}

func user(w http.ResponseWriter, r *http.Request) {

	var errMessage ErrMessage
	var user User

	switch r.Method {
	case "GET":

		dataCTX := r.Context().Value(ContextUserKey)

		if userBeta, ok := dataCTX.(User); !ok {
			errMessage.Message = http.StatusText(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(errMessage)
			return
		} else {
			user = userBeta
		}

		err := json.NewEncoder(w).Encode(user)

		if err != nil {
			errMessage.Message = http.StatusText(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(errMessage)
			return
		}

	}

}

func login(w http.ResponseWriter, r *http.Request) {

	var errMessage ErrMessage

	TokenValue := struct {
		Token string
	}{}

	switch r.Method {
	case "POST":

		var user User

		user.Username = r.Header.Get("username")
		user.Password = r.Header.Get("password")

		row := db.QueryRow("SELECT users.id,users.role FROM users WHERE users.username=? AND users.password=?", user.Username, user.Password)

		err := row.Scan(&user.ID, &user.Role)

		if err != nil {
			errMessage.Message = http.StatusText(http.StatusBadRequest)
			json.NewEncoder(w).Encode(errMessage)
			return
		}

		user.Token, err = generateToken(user.ID, user.Username, user.Role)

		if err != nil {
			errMessage.Message = http.StatusText(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(errMessage)
			return
		}

		TokenValue.Token = user.Token

		err = json.NewEncoder(w).Encode(TokenValue)

		if err != nil {
			errMessage.Message = http.StatusText(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(errMessage)
			return
		}
	}

}

func register(w http.ResponseWriter, r *http.Request) {

	var errMessage ErrMessage

	TokenValue := struct {
		Token string
	}{}

	switch r.Method {
	case "POST":

		var user User

		user.Username = r.Header.Get("username")
		user.Password = r.Header.Get("password")

		err := addUser(&user)

		if err != nil {
			errMessage.Message = err.Error()
			json.NewEncoder(w).Encode(errMessage)
			return
		}

		user.Token, err = generateToken(user.ID, user.Username, user.Role)

		if err != nil {
			errMessage.Message = http.StatusText(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(errMessage)
			return
		}

		TokenValue.Token = user.Token

		err = json.NewEncoder(w).Encode(TokenValue)

		if err != nil {
			errMessage.Message = http.StatusText(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(errMessage)
			return
		}

	}

}

func logout(w http.ResponseWriter, r *http.Request) {

	var errMessage ErrMessage

	Mensaje := struct {
		Mensaje string
	}{}

	var user User

	switch r.Method {
	case "GET":

		dataCTX := r.Context().Value(ContextUserKey)

		if userBeta, ok := dataCTX.(User); !ok {
			errMessage.Message = http.StatusText(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(errMessage)
			return
		} else {
			user = userBeta
		}

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

		Mensaje.Mensaje = "Sesión cerrada"

		err = json.NewEncoder(w).Encode(Mensaje)

		if err != nil {
			errMessage.Message = http.StatusText(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(errMessage)
			return
		}

	}

}

func middlewareGetUser(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		ErrMensaje := struct {
			Mensaje string
		}{}

		token := r.Header.Get("Authorization header")

		check := checkIfTokenIsInBlackList(token)

		if !check {
			ErrMensaje.Mensaje = "El token no es válido"
			json.NewEncoder(w).Encode(ErrMensaje)
			return
		}

		user, err := extractUserFromClaims(token)

		if err != nil {
			ErrMensaje.Mensaje = err.Error()
			json.NewEncoder(w).Encode(ErrMensaje)
			return
		}

		deadline, err := time.Parse(time.ANSIC, user.Deadline)

		if err != nil {
			ErrMensaje.Mensaje = http.StatusText(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(ErrMensaje)
			return
		}

		checkTime := time.Now().Local().After(deadline)

		if !checkTime {
			ErrMensaje.Mensaje = "El token no es válido"
			json.NewEncoder(w).Encode(ErrMensaje)
			return
		}

		user.Token = token

		ctx := context.WithValue(r.Context(), ContextUserKey, user)

		next.ServeHTTP(w, r.WithContext(ctx))

	})

}

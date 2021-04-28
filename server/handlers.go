package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/cfabrica46/chat-web-socket/server/database"
	"github.com/cfabrica46/chat-web-socket/server/middlewares"
	"github.com/cfabrica46/chat-web-socket/server/token"
)

type ErrMessage struct {
	Message string
}

func user(w http.ResponseWriter, r *http.Request) {

	log.SetFlags(log.Lshortfile)

	var errMessage ErrMessage
	var user database.User

	switch r.Method {
	case "GET":

		dataCTX := r.Context().Value(middlewares.ContextUserKey)

		if userBeta, ok := dataCTX.(database.User); !ok {
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

		var user database.User

		user.Username = r.Header.Get("username")
		user.Password = r.Header.Get("password")

		row := db.D.QueryRow("SELECT users.id,users.role FROM users WHERE users.username=? AND users.password=?", user.Username, user.Password)

		err := row.Scan(&user.ID, &user.Role)

		if err != nil {
			errMessage.Message = http.StatusText(http.StatusBadRequest)
			json.NewEncoder(w).Encode(errMessage)
			return
		}

		user.Token, err = token.GenerateToken(user.ID, user.Username, user.Role)

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

		var user database.User

		user.Username = r.Header.Get("username")
		user.Password = r.Header.Get("password")

		err := database.AddUser(&user, db.D)

		if err != nil {
			errMessage.Message = err.Error()
			json.NewEncoder(w).Encode(errMessage)
			return
		}

		user.Token, err = token.GenerateToken(user.ID, user.Username, user.Role)

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

	var user database.User

	switch r.Method {
	case "GET":

		dataCTX := r.Context().Value(middlewares.ContextUserKey)

		if userBeta, ok := dataCTX.(database.User); !ok {
			errMessage.Message = http.StatusText(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(errMessage)
			return
		} else {
			user = userBeta
		}

		stmt, err := db.D.Prepare("INSERT INTO black_list(token) VALUES (?)")

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

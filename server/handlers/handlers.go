package handlers

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

type Message struct {
	Message string
}

func User(w http.ResponseWriter, r *http.Request) {

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

func Login(w http.ResponseWriter, r *http.Request) {

	var errMessage ErrMessage

	TokenValue := struct {
		Token string
	}{}

	switch r.Method {
	case "POST":

		var user database.User

		var err error

		dataCTX := r.Context().Value(middlewares.ContextUserKey)

		if userBeta, ok := dataCTX.(database.User); !ok {
			errMessage.Message = http.StatusText(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(errMessage)
			return

		} else {
			user = userBeta
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

func Register(w http.ResponseWriter, r *http.Request) {

	var errMessage ErrMessage

	TokenValue := struct {
		Token string
	}{}

	switch r.Method {
	case "POST":

		var user database.User
		var err error

		dataCTX := r.Context().Value(middlewares.ContextUserKey)

		if userBeta, ok := dataCTX.(database.User); !ok {

			errMessage.Message = http.StatusText(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(errMessage)
			return

		} else {
			user = userBeta
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

func Logout(w http.ResponseWriter, r *http.Request) {

	var errMessage ErrMessage

	var message Message

	switch r.Method {
	case "GET":

		message.Message = "Sesión Cerrada"

		err := json.NewEncoder(w).Encode(message)

		if err != nil {
			errMessage.Message = http.StatusText(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(errMessage)
			return
		}

	}

}

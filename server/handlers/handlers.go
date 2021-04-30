package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/cfabrica46/chat-web-socket/server/database"
	"github.com/cfabrica46/chat-web-socket/server/middlewares"
	"github.com/cfabrica46/chat-web-socket/server/token"
)

type ErrMessage struct {
	Message string
}

func User(w http.ResponseWriter, r *http.Request) {

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

		row := database.DB.QueryRow("SELECT users.id,users.role FROM users WHERE users.username=? AND users.password=?", user.Username, user.Password)

		err = row.Scan(&user.ID, &user.Role)

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

		err = database.AddUser(&user, database.DB)

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

func Logout(w http.ResponseWriter, r *http.Request) {

	var errMessage ErrMessage
	var user database.User

	Message := struct {
		Message string
	}{}

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

		stmt, err := database.DB.Prepare("INSERT INTO black_list(token) VALUES (?)")

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

		Message.Message = "Sesión Cerrada"

		err = json.NewEncoder(w).Encode(Message)

		if err != nil {
			errMessage.Message = http.StatusText(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(errMessage)
			return
		}

	}

}

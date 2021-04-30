package chat

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/cfabrica46/chat-web-socket/server/database"
	"github.com/cfabrica46/chat-web-socket/server/middlewares"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

type ErrMessage struct {
	Message string
}

type Conex struct {
	User database.User
	C    *websocket.Conn
}

var rooms = make(map[int][]Conex)

var upgrader = websocket.Upgrader{}

func Chat(w http.ResponseWriter, r *http.Request) {

	var errMessage ErrMessage
	var user database.User

	vars := mux.Vars(r)

	idString := vars["id"]

	idRoom, err := strconv.Atoi(idString)

	if err != nil {
		errMessage.Message = http.StatusText(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(errMessage)
		return
	}

	dataCTX := r.Context().Value(middlewares.ContextUserKey)

	if userBeta, ok := dataCTX.(database.User); !ok {

		errMessage.Message = http.StatusText(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(errMessage)
		return

	} else {
		user = userBeta
	}

	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		errMessage.Message = http.StatusText(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(errMessage)
		return

	}

	defer conn.Close()

	conex := Conex{
		User: user,
		C:    conn,
	}

	if len(rooms[idRoom]) >= 2 {

		errMessage.Message = "La sala ya esta llena"
		json.NewEncoder(w).Encode(errMessage)
		return

	}

	rooms[idRoom] = append(rooms[idRoom], conex)

	SendMessageConect(conex, idRoom)

	for {
		_, message, err := conn.ReadMessage()

		if err != nil {
			Disconect(conex, idRoom)
			SendMessageDisconect(conex, idRoom)
			return
		}

		err = SendMessage(conex, idRoom, message)

		if err != nil {
			errMessage.Message = "Hubo un error al enviar el mensaje"
			json.NewEncoder(w).Encode(errMessage)
		}

	}

}

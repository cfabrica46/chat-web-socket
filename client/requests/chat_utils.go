package requests

import (
	"encoding/json"
	"fmt"

	"github.com/gorilla/websocket"
)

func reciveMessage(ws *websocket.Conn, c chan bool) {

	var finish bool

	Message := struct {
		Message string
	}{}

	for !finish {
		select {
		case <-c:
			finish = true
			return
		default:

		}

		_, message, err := ws.ReadMessage()

		if err != nil {
			fmt.Print("")
		} else {

			err = json.Unmarshal(message, &Message)

			if err != nil {
				fmt.Println("\rError al recivir mensaje")
				fmt.Print("\r> ")
			}

			fmt.Printf("\r%s\n", Message.Message)
			fmt.Print("\r> ")

		}

	}

}

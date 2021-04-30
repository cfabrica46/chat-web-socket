package requests

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
)

type User struct {
	ID                                        int
	Username, Password, Deadline, Role, Token string
}

type ErrMessage struct {
	Message string
}

func Login(username, password, url string) (token string, err error) {

	fmt.Println()

	var errMessage ErrMessage

	TokenValue := struct {
		Token string
	}{}

	client := &http.Client{
		Timeout: time.Second * 20,
	}

	req, err := http.NewRequest("POST", url, nil)

	if err != nil {
		return
	}

	req.Header.Set("username", username)
	req.Header.Set("password", password)

	resp, err := client.Do(req)

	if err != nil {
		return
	}

	defer resp.Body.Close()

	dataJSON, err := io.ReadAll(resp.Body)

	if err != nil {
		return
	}

	err = json.Unmarshal(dataJSON, &TokenValue)

	if err != nil {
		return
	}

	if TokenValue.Token == "" {

		err = json.Unmarshal(dataJSON, &errMessage)

		if err != nil {
			return
		}

		fmt.Println(errMessage.Message)

	}

	token = TokenValue.Token

	return
}

func Profile(token string) (err error) {

	var user User

	var errMessage ErrMessage

	client := &http.Client{
		Timeout: time.Second * 20,
	}

	req, err := http.NewRequest("GET", "http://localhost:8080/user", nil)

	if err != nil {
		return
	}

	req.Header.Set("Authorization-header", token)

	resp, err := client.Do(req)

	if err != nil {
		return
	}

	defer resp.Body.Close()

	dataJSON, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return
	}

	err = json.Unmarshal(dataJSON, &user)

	if err != nil {
		log.Fatal(err)
		return
	}

	if user.ID == 0 && user.Username == "" && user.Password == "" && user.Token == "" {

		err = json.Unmarshal(dataJSON, &errMessage)

		if err != nil {
			return
		}

		fmt.Println(errMessage.Message)

		return

	}

	fmt.Printf("Bienvenido %s %s tu ID es: %d y tu Token es: %s\n", user.Role, user.Username, user.ID, user.Token)

	return

}

func LogOut(token string) (err error) {

	var errMessage ErrMessage

	Message := struct {
		Message string
	}{}

	client := &http.Client{
		Timeout: time.Second * 20,
	}

	req, err := http.NewRequest("GET", "http://localhost:8080/logout", nil)

	if err != nil {
		return
	}

	req.Header.Set("Authorization-header", token)

	resp, err := client.Do(req)

	if err != nil {
		return
	}

	defer resp.Body.Close()

	dataJSON, err := io.ReadAll(resp.Body)

	if err != nil {
		log.Fatal(err)
		return
	}

	err = json.Unmarshal(dataJSON, &Message)

	if err != nil {
		log.Fatal(err)

		return
	}

	if Message.Message != "" {
		err = json.Unmarshal(dataJSON, &errMessage)

		if err != nil {
			return
		}

		fmt.Println(errMessage.Message)
		return
	}

	fmt.Println(Message.Message)

	return
}

func Chat(token string, idRoom int) (err error) {

	header := make(map[string][]string)

	header["Authorization-header"] = []string{token}

	c := make(chan bool)

	id := strconv.Itoa(idRoom)

	path := fmt.Sprintf("/chat/%s", id)

	u := url.URL{
		Scheme: "ws",
		Host:   ":8080",
		Path:   path,
	}

	ws, _, _ := websocket.DefaultDialer.Dial(u.String(), header)

	defer ws.Close()

	fmt.Println()
	fmt.Println("Has ingresado")
	fmt.Println()

	go reciveMessage(ws, c)

	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		fmt.Print("\r> ")

		var dataJSON []byte

		Message := struct {
			Message string
		}{}

		Message.Message = scanner.Text()

		dataJSON, err = json.Marshal(Message)

		if err != nil {
			return
		}

		err = ws.WriteMessage(websocket.TextMessage, dataJSON)

		if err != nil {
			return
		}

		if scanner.Text() == "--salir--" {
			ws.Close()
			c <- true
			return
		}
	}

	return
}

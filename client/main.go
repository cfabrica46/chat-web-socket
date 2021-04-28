package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type User struct {
	ID                                        int
	Username, Password, Deadline, Role, Token string
}

type ErrMessage struct {
	Message string
}

func main() {

	log.SetFlags(log.Lshortfile)

	var election int
	var exit bool

	for !exit {

		fmt.Println()
		fmt.Println("Bienvenido")
		fmt.Println("¿Qué deseas hacer?")
		fmt.Println()

		fmt.Println("1.Iniciar Sesión")
		fmt.Println("2.Registrarse")
		fmt.Println("0.Salir")
		fmt.Println()

		fmt.Print("> ")

		fmt.Scan(&election)

		fmt.Println()

		switch election {

		case 0:

			exit = true

		case 1:

			var username, password string

			fmt.Printf("\nINGRESA TUS DATOS\n")
			fmt.Printf("Username: ")
			fmt.Scan(&username)
			fmt.Printf("Password: ")
			fmt.Scan(&password)

			token, err := login(username, password, "http://localhost:8080/login")

			if err != nil {
				log.Fatal(err)
			}

			if token == "" {
				break
			}

			err = profile(token)

			if err != nil {
				log.Fatal(err)
			}

			for !exit {
				err = loopIntoProfile(token, &exit)

				if err != nil {
					log.Fatal(err)
				}
			}

		case 2:

			var username, password string

			fmt.Printf("\nINGRESA TUS DATOS\n")
			fmt.Printf("Username: ")
			fmt.Scan(&username)
			fmt.Printf("Password: ")
			fmt.Scan(&password)

			token, err := login(username, password, "http://localhost:8080/register")

			if err != nil {
				log.Fatal(err)
			}

			if token == "" {
				break
			}

			err = profile(token)

			if err != nil {
				log.Fatal(err)
			}

			for !exit {
				err = loopIntoProfile(token, &exit)

				if err != nil {
					log.Fatal(err)
				}

			}

		default:

			fmt.Println("Seleccione una opción válida")

		}
	}

}

func loopIntoProfile(token string, exit *bool) (err error) {

	var election int

	fmt.Println()
	fmt.Println("¿Qué deseas hacer?")
	fmt.Println()

	fmt.Println("1.Cerrar Sesión")
	fmt.Println("0.Salir")
	fmt.Println()

	fmt.Print("> ")

	fmt.Scan(&election)

	fmt.Println()

	switch election {
	case 0:

		*exit = true

	case 1:

		err = cerrarSesión(token)

		if err != nil {
			return
		}

		*exit = true

	default:

		fmt.Println("Seleccione una opción válida")

	}

	return
}

func login(username, password, url string) (token string, err error) {

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

func profile(token string) (err error) {

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

	fmt.Printf("%s\n", dataJSON)

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

func cerrarSesión(token string) (err error) {

	var errMessage ErrMessage

	Mensaje := struct {
		Mensaje string
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

	err = json.Unmarshal(dataJSON, &Mensaje)

	if err != nil {
		log.Fatal(err)

		return
	}

	if Mensaje.Mensaje != "" {
		err = json.Unmarshal(dataJSON, &errMessage)

		if err != nil {
			return
		}

		fmt.Println(errMessage.Message)
		return
	}

	fmt.Println(Mensaje.Mensaje)

	return
}

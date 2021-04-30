package main

import (
	"fmt"
	"log"

	"github.com/cfabrica46/chat-web-socket/client/requests"
)

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

			token, err := requests.Login(username, password, "http://localhost:8080/login")

			if err != nil {
				log.Fatal(err)
			}

			if token == "" {
				break
			}

			err = requests.Profile(token)

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

			token, err := requests.Login(username, password, "http://localhost:8080/register")

			if err != nil {
				log.Fatal(err)
			}

			if token == "" {
				break
			}

			err = requests.Profile(token)

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

	fmt.Println("1.Chatear")
	fmt.Println("2.Cerrar Sesión")
	fmt.Println("0.Salir")
	fmt.Println()

	fmt.Print("> ")

	fmt.Scan(&election)

	fmt.Println()

	switch election {
	case 0:

		*exit = true

	case 1:

		var idRoom int

		fmt.Println()

		fmt.Println("¿A que sala quieres ingresar?")
		fmt.Println()
		fmt.Print("> ")

		fmt.Scan(&idRoom)

		err = requests.Chat(token, idRoom)

		if err != nil {
			return
		}

	case 2:

		err = requests.LogOut(token)

		if err != nil {
			return
		}

		*exit = true

	default:

		fmt.Println("Seleccione una opción válida")

	}

	return
}

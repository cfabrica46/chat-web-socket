package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/cfabrica46/chat-web-socket/server/database"
	"github.com/cfabrica46/chat-web-socket/server/middlewares"
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

	muxUser := http.HandlerFunc(user)
	muxLogin := http.HandlerFunc(login)
	muxRegister := http.HandlerFunc(register)
	muxLogout := http.HandlerFunc(logout)

	r.Handle("/user", middlewares.GetUser(muxUser, db.D))
	r.Handle("/login", middlewares.LoginPassword(muxLogin, db.D))
	r.Handle("/register", middlewares.LoginPassword(muxRegister, db.D))
	r.Handle("/logout", middlewares.GetUser(muxLogout, db.D))

	fmt.Println("Listening on localhost:8080")

	http.ListenAndServe(":8080", r)
}

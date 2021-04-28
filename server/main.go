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

	r.HandleFunc("/login", login)
	r.HandleFunc("/register", register)

	muxUser := http.HandlerFunc(user)
	muxLogout := http.HandlerFunc(logout)

	r.Handle("/user", middlewares.GetUser(muxUser, db.D))
	r.Handle("/logout", middlewares.GetUser(muxLogout, db.D))

	fmt.Println("Listening on localhost:8080")

	http.ListenAndServe(":8080", r)
}

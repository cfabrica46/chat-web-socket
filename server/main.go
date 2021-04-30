package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/cfabrica46/chat-web-socket/server/chat"
	"github.com/cfabrica46/chat-web-socket/server/database"
	"github.com/cfabrica46/chat-web-socket/server/handlers"
	"github.com/cfabrica46/chat-web-socket/server/middlewares"
	"github.com/gorilla/mux"
)

var db database.DB

func main() {

	log.SetFlags(log.Lshortfile)

	d, err := database.Open()

	if err != nil {
		log.Fatal(err)
	}

	db = database.DB{
		D: d,
	}

	go database.CleanBlackList(db.D)

	r := mux.NewRouter()

	mUser := r.PathPrefix("/").Subrouter()

	mUser.HandleFunc("/user", handlers.User)
	mUser.HandleFunc("/logout", handlers.Logout)
	mUser.HandleFunc("/chat/{id:[0-9]+}", chat.Chat)
	mUser.Use(middlewares.GetUser(db.D))

	mLogin := r.PathPrefix("/").Subrouter()

	mLogin.HandleFunc("/login", handlers.Login)
	mLogin.HandleFunc("/register", handlers.Register)
	mLogin.Use(middlewares.LoginPassword(db.D))

	fmt.Println("Listening on localhost:8080")

	http.ListenAndServe(":8080", r)
}

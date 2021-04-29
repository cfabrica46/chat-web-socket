package main

import (
	"fmt"
	"log"
	"net/http"

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

	s := r.PathPrefix("/").Subrouter()

	s.HandleFunc("/user", handlers.User)
	s.HandleFunc("/logout", handlers.Logout)
	s.Use(middlewares.GetUser(db.D))

	muxLogin := http.HandlerFunc(handlers.Login)
	muxRegister := http.HandlerFunc(handlers.Register)

	r.Handle("/login", middlewares.LoginPassword(muxLogin, db.D))
	r.Handle("/register", middlewares.LoginPassword(muxRegister, db.D))

	//muxRoom := http.HandlerFunc(handlers.Room)

	//r.Handle("/room/{id:[0-9]+}", middlewares.GetUser(muxRoom, db.D))

	fmt.Println("Listening on localhost:8080")

	http.ListenAndServe(":8080", r)
}

package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type User struct {
	ID                        int
	Username, Password, Token string
}

type ErrMessage struct {
	Message string
}

var db *sql.DB

func main() {

	var err error

	r := mux.NewRouter()

	db, err = open()

	if err != nil {
		log.Fatal(err)
	}

	go cleanBlackList()

	r.HandleFunc("/login", login)
}

func user(w http.ResponseWriter, r *http.Request) {

	var errMessage ErrMessage

	switch r.Method {
	case "GET":

	}

}

func login(w http.ResponseWriter, r *http.Request) {

	var errMessage ErrMessage

	TokenValue := struct {
		Token string
	}{}

	switch r.Method {
	case "POST":

	}

}

func register(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case "POST":
	}

}

func logout(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case "GET":
	}

}

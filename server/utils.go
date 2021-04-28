package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
)

func open() (databases *sql.DB, err error) {
	archivo, err := os.Open("databases.db")

	if err != nil {
		if os.IsNotExist(err) {

			databases, err = migracion()

			if err != nil {

				archivo.Close()
				return
			}

			return
		}
		return
	}
	defer archivo.Close()

	databases, err = sql.Open("sqlite3", "./databases.db?_foreign_keys=on")

	if err != nil {
		return
	}

	return
}

func migracion() (databases *sql.DB, err error) {

	archivoDB, err := os.Create("databases.db")

	if err != nil {
		return
	}
	archivoDB.Close()

	databases, err = sql.Open("sqlite3", "./databases.db?_foreign_keys=on")

	if err != nil {
		return
	}

	dataSQL, err := ioutil.ReadFile("databases.sql")

	if err != nil {
		return
	}

	_, err = databases.Exec(string(dataSQL))

	if err != nil {
		return
	}

	return

}

func generateToken(id int, username, role string) (tokenString string, err error) {

	secret, err := ioutil.ReadFile("key.pem")

	if err != nil {
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":        id,
		"username":  username,
		"dead-line": time.Now().Add(1 * time.Hour).Format(time.ANSIC),
		"role":      role,
		"uuid":      uuid.NewString(),
	})

	tokenString, err = token.SignedString(secret)

	if err != nil {
		return
	}

	return

}

func extractUserFromClaims(tokenString string) (user User, err error) {

	token, err := jwt.Parse(tokenString, keyFunc)

	if err != nil {
		return
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {

		if idFloat64, ok := claims["id"].(float64); ok {

			user.ID = int(idFloat64)

		} else {

			err = fmt.Errorf("error en obtener claims")
			return
		}

		if username, ok := claims["username"].(string); ok {

			user.Username = username

		} else {

			err = fmt.Errorf("error en obtener claims")
			return
		}

		if deadline, ok := claims["dead-line"].(string); ok {

			user.Deadline = deadline

		} else {

			err = fmt.Errorf("error en obtener claims")
			return
		}

		if role, ok := claims["role"].(string); ok {

			user.Role = role

		} else {

			err = fmt.Errorf("error en obtener claims")
			return
		}

	} else {

		err = fmt.Errorf("error en obtener claims")
	}

	return
}

func checkIfTokenIsInBlackList(token string) (check bool) {

	var tokenAux string

	row := db.QueryRow("SELECT token FROM black_list WHERE token = ?", token)

	err := row.Scan(&tokenAux)

	if err != nil {
		if err != sql.ErrNoRows {
			return
		}
		check = true
	}

	return
}

func addUser(user *User) (err error) {

	user.Role = "member"

	stmt, err := db.Prepare("INSERT INTO users(username,password,role) VALUES (?,?,?)")

	if err != nil {
		return
	}

	result, err := stmt.Exec(user.Username, user.Password, user.Role)

	if err != nil {
		return
	}

	id64, err := result.LastInsertId()

	if err != nil {
		return
	}

	user.ID = int(id64)

	return
}

func cleanBlackList() {

	for {

		stmt, err := db.Prepare("DELETE FROM black_list")

		if err != nil {
			log.Fatal(err)
		}

		_, err = stmt.Exec()

		if err != nil {
			log.Fatal(err)
		}

		time.Sleep(time.Hour * 1)

	}

}

func keyFunc(token *jwt.Token) (interface{}, error) {

	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {

		return nil, fmt.Errorf("inesperado metodo : %v", token.Header["alg"])

	}

	secret, err := ioutil.ReadFile("key.pem")

	if err != nil {
		return nil, err
	}

	return secret, nil
}

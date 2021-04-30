package database

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type User struct {
	ID                                        int
	Username, Password, Deadline, Role, Token string
}

func CheckIfTokenIsInBlackList(token string, db *sql.DB) (check bool) {

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

func AddUser(user *User, db *sql.DB) (err error) {

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

func CleanBlackList(db *sql.DB) {

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

func CheckIfRoomExist(db *sql.DB, id int) (check bool, err error) {

	row := db.QueryRow("SELECT id FROM rooms WHERE id = ?", id)

	err = row.Scan(&id)

	if err != nil {
		if err == sql.ErrNoRows {

			err = nil
			return
		}

		return
	}

	check = true

	return
}

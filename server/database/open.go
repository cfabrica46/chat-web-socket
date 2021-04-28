package database

import (
	"database/sql"
	"io/ioutil"
	"os"
)

func Open() (databases *sql.DB, err error) {
	archivo, err := os.Open("databases.db")

	if err != nil {
		if os.IsNotExist(err) {

			databases, err = Migracion()

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

func Migracion() (databases *sql.DB, err error) {

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

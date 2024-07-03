package config

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

const JwtSecret = "your_secret_key"

func InitDatabase() {
	var err error
	DB, err = sql.Open("sqlite3", "./users.db")
	if err != nil {
		log.Fatal(err)
	}

	createUserTableSQL := `CREATE TABLE IF NOT EXISTS users (
        "id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
        "username" TEXT NOT NULL UNIQUE,
        "password" TEXT NOT NULL
    );`

	statement, err := DB.Prepare(createUserTableSQL)
	if err != nil {
		log.Fatal(err.Error())
	}
	statement.Exec()
}

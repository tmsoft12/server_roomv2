package database

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

// Veritabanı bağlantısını başlatan fonksiyon

// Veritabanı bağlantısını başlatan fonksiyon
func InitDatabase() error {
	var err error
	db, err = sql.Open("sqlite3", "./server.sql")
	if err != nil {
		return err
	}

	// Veritabanı bağlantısını doğrulamak için ping atıyoruz
	if err = db.Ping(); err != nil {
		return err
	}
	return nil
}

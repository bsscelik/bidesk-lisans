package main

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/lib/pq"
)

var db *sql.DB

func InitDB() {
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	var err error

	// 🔥 retry loop
	for i := 0; i < 10; i++ {
		db, err = sql.Open("postgres", connStr)
		if err == nil {
			err = db.Ping()
			if err == nil {
				break
			}
		}

		fmt.Println("DB bekleniyor...")
		time.Sleep(2 * time.Second)
	}

	if err != nil {
		panic("DB baglantisi kurulamadı: " + err.Error())
	}

	createTable()
}

func createTable() {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS licenses (
			code TEXT PRIMARY KEY,
			start_at TIMESTAMP,
			end_at TIMESTAMP,
			period TEXT,
			active BOOLEAN,
			db_name TEXT
		)`,
		`CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			username TEXT UNIQUE NOT NULL,
			password_hash TEXT NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS sessions (
			id TEXT PRIMARY KEY,
			user_id INTEGER REFERENCES users(id),
			expires_at TIMESTAMP NOT NULL
		)`,
	}
	for _, query := range queries {
		_, err := db.Exec(query)
		if err != nil {
			panic(err)
		}
	}
}
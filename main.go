package main

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

func main() {
	InitDB()

	http.HandleFunc("/register", register)
	http.HandleFunc("/login", login)

	http.HandleFunc("/license", authMiddleware(createLicense))
	http.HandleFunc("/license/get", authMiddleware(getLicense))

	log.Println("running :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func GenerateDBName(code string) string {
	b := make([]byte, 4)
	rand.Read(b)
	return fmt.Sprintf("t_%s_%s", code, hex.EncodeToString(b))
}

func calculateEnd(start time.Time, period string) time.Time {
	switch period {
	case "monthly":
		return start.AddDate(0, 1, 0)
	case "yearly":
		return start.AddDate(1, 0, 0)
	default:
		return start
	}
}

func createLicense(w http.ResponseWriter, r *http.Request) {
	var l License
	json.NewDecoder(r.Body).Decode(&l)

	l.StartAt = time.Now()
	l.EndAt = calculateEnd(l.StartAt, l.Period)
	l.DBName = GenerateDBName(l.Code)
	l.Active = true

	_, err := db.Exec(
		"INSERT INTO licenses (code, start_at, end_at, period, active, db_name) VALUES ($1,$2,$3,$4,$5,$6)",
		l.Code, l.StartAt, l.EndAt, l.Period, l.Active, l.DBName,
	)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// 🔥 DB create (template1 kullanıyoruz basitlik için)
	_, err = db.Exec(fmt.Sprintf("CREATE DATABASE %s TEMPLATE template1", l.DBName))
	if err != nil {
		log.Println("db create error:", err)
	}

	json.NewEncoder(w).Encode(l)
}

func getLicense(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")

	var l License
	err := db.QueryRow(
		"SELECT code, start_at, end_at, period, active, db_name FROM licenses WHERE code=$1",
		code,
	).Scan(&l.Code, &l.StartAt, &l.EndAt, &l.Period, &l.Active, &l.DBName)

	if err == sql.ErrNoRows {
		http.NotFound(w, r)
		return
	}

	json.NewEncoder(w).Encode(l)
}
package main

import "time"

type License struct {
	Code    string    `json:"code"`
	StartAt time.Time `json:"start_at"`
	EndAt   time.Time `json:"end_at"`
	Period  string    `json:"period"`
	Active  bool      `json:"active"`
	DBName  string    `json:"db_name"`
}
package data

import "database/sql"

type College struct {
	ID     int64  `json:"id"`
	Name   string `json:"name"`
	Domain string `json:"domain"`
}

type CollegeModel struct {
	DB *sql.DB
}

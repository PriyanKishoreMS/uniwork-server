package data

import (
	"database/sql"
	"time"
)

type Record struct {
	ID          int64     `json:"id"`
	UserID      int64     `json:"user_id"`
	CollegeID   int64     `json:"college_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Category    string    `json:"category"`
	Price       int64     `json:"price"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"time"`
	Expiry      string    `json:"expiry"`
	Images      []string  `json:"images,omitempty"`
}

type ServiceModel struct {
	DB *sql.DB
}

package data

import (
	"database/sql"
	"time"
)

type User struct {
	ID         int64     `json:"id"`
	CollegeID  int64     `json:"college_id"`
	Name       string    `json:"name"`
	Email      string    `json:"email"`
	ProfilePic string    `json:"profile_pic,omitempty"`
	Dept       string    `json:"dept"`
	Review     float64   `json:"review"`
	CreatedAt  time.Time `json:"created_at"`
}

type UserModel struct {
	DB *sql.DB
}

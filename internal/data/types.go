package data

import (
	"encoding/json"
	"time"
)

type User struct {
	ID             int64     `json:"id"`
	CollegeID      int64     `json:"college_id"`
	CollegeName    string    `json:"college_name"`
	Name           string    `json:"name" validate:"required"`
	Email          string    `json:"email,omitempty" validate:"required,email"`
	Mobile         string    `json:"mobile,omitempty"`
	Avatar         string    `json:"avatar,omitempty"`
	Dept           string    `json:"dept" validate:"required"`
	TasksCompleted int       `json:"tasks_completed"`
	Earned         int64     `json:"earned"`
	Rating         float64   `json:"rating"`
	Status         string    `json:"status"`
	CreatedAt      time.Time `json:"created_at,omitempty"`
	Version        int       `json:"version,omitempty"`
}

// struct to add the task to db
type Task struct {
	ID          int64     `form:"id"`
	UserID      int64     `form:"user_id" validate:"required"`
	CollegeID   int64     `form:"college_id" validate:"required"`
	Title       string    `form:"title" validate:"required"`
	Description string    `form:"description"`
	Category    string    `form:"category" validate:"required"`
	Price       int64     `form:"price" validate:"required"`
	Status      string    `form:"status" validate:"required"`
	CreatedAt   time.Time `form:"time"`
	Expiry      string    `form:"expiry" validate:"required"`
	Images      []string  `form:"images,omitempty"`
	Files       []string  `form:"files,omitempty"`
}

//	type requester struct {
//		RequestStatus   string `json:"request_status"`
//		RequesterName   string `json:"requester_name"`
//		RequesterAvatar string `json:"requester_avatar"`
//	}

// viewing details of single task
type GetTaskResponse struct {
	ID          int64           `json:"id"`
	UserID      int64           `json:"user_id"`
	CollegeName string          `json:"college_name"`
	Title       string          `json:"title"`
	Description string          `json:"description"`
	Category    string          `json:"category"`
	Price       int64           `json:"price"`
	Status      string          `json:"status"`
	CreatedAt   time.Time       `json:"time"`
	Expiry      time.Time       `json:"expiry"`
	Images      []string        `json:"images,omitempty"`
	Files       []string        `json:"files,omitempty"`
	UserName    string          `json:"user_name"`
	UserAvatar  string          `json:"user_avatar"`
	UserRating  float64         `json:"user_rating"`
	Requesters  json.RawMessage `json:"requesters,omitempty"`
}

// listing all the tasks in the home page
type TaskWithUser struct {
	ID          int64     `json:"id"`
	UserID      int64     `json:"user_id" validate:"required"`
	CollegeID   int64     `json:"college_id" validate:"required"`
	Title       string    `json:"title" validate:"required"`
	Description string    `json:"description"`
	Category    string    `json:"category" validate:"required"`
	Price       int64     `json:"price" validate:"required"`
	Status      string    `json:"status" validate:"required"`
	CreatedAt   time.Time `json:"time"`
	Expiry      time.Time `json:"expiry" validate:"required"`
	Images      []string  `json:"images,omitempty"`
	UserName    string    `json:"user_name"`
	UserAvatar  string    `json:"user_avatar"`
	UserRating  float64   `json:"user_rating"`
}

type Status string

const (
	StatusPending  Status = "pending"
	StatusAccepted Status = "accepted"
	StatusRejected Status = "rejected"
)

type TaskRequest struct {
	ID     int64  `json:"id"`
	TaskID int64  `json:"task_id" validate:"required"`
	UserID int64  `json:"user_id" validate:"required"`
	Status Status `json:"status"`
}

type OrderData struct {
	OrderID       string
	TaskOwnerID   int64
	Amount        int64
	TaskRequestID int64
}

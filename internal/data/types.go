package data

import "time"

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
	Expiry      time.Time `form:"expiry" validate:"required"`
	Images      []string  `form:"images,omitempty"`
	Files       []string  `form:"files,omitempty"`
}

// viewing details of single task
type GetTaskResponse struct {
	ID          int64     `json:"id"`
	UserID      int64     `json:"user_id"`
	CollegeName string    `json:"college_name"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Category    string    `json:"category"`
	Price       int64     `json:"price"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"time"`
	Expiry      time.Time `json:"expiry"`
	Images      []string  `json:"images,omitempty"`
	Files       []string  `json:"files,omitempty"`
	UserName    string    `json:"user_name"`
	UserAvatar  string    `json:"user_avatar"`
	UserRating  float64   `json:"user_rating"`
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

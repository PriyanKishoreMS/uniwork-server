package data

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrEditConflict   = errors.New("edit conflict")
)

var TaskCategories []string = []string{
	"Academic Assistance",
	"Tutor Home/Virtual",
	"Books Rent/Buy",
	"Vehicle Rent",
	"Document Printing",
	"Resume Creation",
	"Job Search Support",
	"Grocery Shopping",
	"Fashion",
	"Social Media",
	"IT Support",
	"Graphic Design",
	"Delivery",
	"Ride sharing",
	"Catering/Cooking",
}

func handlectx() (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	return ctx, cancel
}

type Models struct {
	Tasks        TaskModel
	TaskRequests TaskRequestModel
	Users        UserModel
	Colleges     CollegeModel
	FcmTokens    FcmModel
}

func NewModel(db *sql.DB) Models {
	return Models{
		Tasks:        TaskModel{DB: db},
		TaskRequests: TaskRequestModel{DB: db},
		Users:        UserModel{DB: db},
		Colleges:     CollegeModel{DB: db},
		FcmTokens:    FcmModel{DB: db},
	}
}

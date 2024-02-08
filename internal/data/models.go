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

func handlectx() (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	return ctx, cancel
}

type Models struct {
	Tasks    TaskModel
	Users    UserModel
	Colleges CollegeModel
}

func NewModel(db *sql.DB) Models {
	return Models{
		Tasks:    TaskModel{DB: db},
		Users:    UserModel{DB: db},
		Colleges: CollegeModel{DB: db},
	}
}

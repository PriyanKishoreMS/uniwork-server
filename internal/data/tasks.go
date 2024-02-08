package data

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/labstack/gommon/log"
)

var TaskCategories []string = []string{
	"Academic Assistance",
	"Tutor Home/Virtual",
	"Books Rent/Buy",
	"Vechicle Rent",
	"Document Printing",
	"Resume Creation",
	"Job Search support",
	"Grocery Shopping",
	"Fashion",
	"Social Media",
	"IT Support",
	"Graphic Design",
	"Delivery",
	"Ride sharing",
	"Catering/Cooking",
}

type Task struct {
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
}

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

type TaskModel struct {
	DB *sql.DB
}

func (t TaskModel) Create(task *Task) (int64, error) {
	query := `
	INSERT INTO tasks (user_id, college_id, title, description, category, price, status, expiry, images)
	VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	imagesCSV := strings.Join(task.Images, ",")
	log.Info(imagesCSV)

	args := []interface{}{
		task.UserID,
		task.CollegeID,
		task.Title,
		task.Description,
		task.Category,
		task.Price,
		task.Status,
		task.Expiry,
		imagesCSV,
	}

	ctx, cancel := handlectx()
	defer cancel()

	res, err := t.DB.ExecContext(ctx, query, args...)
	if err != nil {
		return -1, err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return -1, err
	}
	log.Info(rowsAffected)

	if rowsAffected == 0 {
		return -1, sql.ErrNoRows
	}

	LastInsertId, err := res.LastInsertId()
	if err != nil {
		return -1, err
	}

	return LastInsertId, nil
}

func (t TaskModel) Get(id int64) (*Task, error) {
	query := `SELECT id, user_id, college_id, title, description, category, price, status, created_at, Expiry, Images
	FROM tasks
	WHERE id=?
	`
	ctx, cancel := handlectx()
	defer cancel()

	var task Task
	var imagesCSV string

	dest := []interface{}{
		&task.ID,
		&task.UserID,
		&task.CollegeID,
		&task.Title,
		&task.Description,
		&task.Category,
		&task.Price,
		&task.Status,
		&task.CreatedAt,
		&task.Expiry,
		&imagesCSV,
	}

	err := t.DB.QueryRowContext(ctx, query, id).Scan(dest...)
	if err != nil {
		return nil, err
	}

	task.Images = strings.Split(imagesCSV, ",")

	return &task, nil
}

func (t TaskModel) Delete(id int64) (int64, error) {
	query := `DELETE FROM tasks WHERE id=?`

	ctx, cancel := handlectx()
	defer cancel()

	res, err := t.DB.ExecContext(ctx, query, id)
	if err != nil {
		return -1, err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return -1, err
	}

	if rowsAffected == 0 {
		return -1, ErrRecordNotFound
	}

	return rowsAffected, nil
}

func (t TaskModel) GetAllInCollege(category string, college_id int64, filters Filters) ([]*TaskWithUser, Metadata, error) {

	query := fmt.Sprint(`
	SELECT 
	COUNT(*) OVER () AS total,
	tasks.id, tasks.college_id, tasks.title, tasks.description, tasks.category, tasks.price, tasks.status, tasks.created_at, tasks.expiry, tasks.images, users.name, users.avatar, users.rating
	FROM tasks
	INNER JOIN users ON users.id=tasks.user_id	
	WHERE tasks.college_id=?
	AND tasks.category IN (` + category + `)
	AND tasks.status='open'
	ORDER BY ` + filters.sortColumn() + " " + filters.sortDirection() + `, id ASC
	LIMIT ? OFFSET ?;
	`)

	ctx, cancel := handlectx()
	defer cancel()

	args := []interface{}{
		college_id,
		filters.limit(),
		filters.offset(),
	}

	rows, err := t.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, Metadata{}, err
	}

	defer rows.Close()

	totalRecords := 0
	tasks := []*TaskWithUser{}
	var imagesCSV string

	for rows.Next() {
		var task TaskWithUser

		err := rows.Scan(
			&totalRecords,
			&task.ID,
			&task.CollegeID,
			&task.Title,
			&task.Description,
			&task.Category,
			&task.Price,
			&task.Status,
			&task.CreatedAt,
			&task.Expiry,
			&imagesCSV,
			&task.UserName,
			&task.UserAvatar,
			&task.UserRating,
		)
		if err != nil {
			return nil, Metadata{}, err
		}

		task.Images = strings.Split(imagesCSV, ",")
		tasks = append(tasks, &task)
	}

	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}
	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)

	return tasks, metadata, nil
}

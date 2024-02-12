package data

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/lib/pq"
)

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

func (t TaskModel) Create(task *Task) error {
	query := `
	INSERT INTO tasks (user_id, college_id, title, description, category, price, status, expiry, images)
	VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9)
	RETURNING id
	`

	args := []interface{}{
		task.UserID,
		task.CollegeID,
		task.Title,
		task.Description,
		task.Category,
		task.Price,
		task.Status,
		task.Expiry,
		pq.Array(task.Images),
	}

	ctx, cancel := handlectx()
	defer cancel()

	return t.DB.QueryRowContext(ctx, query, args...).Scan(&task.ID)

}

func (t TaskModel) Get(id int64) (*Task, error) {
	query := `SELECT id, user_id, college_id, title, description, category, price, status, created_at, Expiry, Images
	FROM tasks
	WHERE id=$1
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

func (t TaskModel) Delete(id int64) error {
	query := `DELETE FROM tasks WHERE id=$1
	RETURNING id
	`

	ctx, cancel := handlectx()
	defer cancel()

	err := t.DB.QueryRowContext(ctx, query, id).Scan(&id)
	if err != nil {
		return err
	}

	return nil

}

func (t TaskModel) GetAllTasksInCollege(category string, college_id int64, filters Filters) ([]*TaskWithUser, Metadata, error) {
	var query string
	if category == "" {
		query = fmt.Sprintf(`
		SELECT
		COUNT(*) OVER () AS total,
		tasks.id, tasks.college_id, tasks.title, tasks.description, tasks.category, tasks.price, tasks.status, tasks.created_at, tasks.expiry, tasks.images, users.name, users.avatar, users.rating
		FROM tasks
		INNER JOIN users ON users.id=tasks.user_id
		WHERE tasks.college_id=$1
		AND tasks.status='open'
		ORDER BY %s %s, id ASC
		LIMIT $2 OFFSET $3;
		`, filters.sortColumn(), filters.sortDirection())
	} else {
		query = fmt.Sprintf(`
	SELECT 
	COUNT(*) OVER () AS total,
	tasks.id, tasks.college_id, tasks.title, tasks.description, tasks.category, tasks.price, tasks.status, tasks.created_at, tasks.expiry, tasks.images, users.name, users.avatar, users.rating
	FROM tasks
	INNER JOIN users ON users.id=tasks.user_id	
	WHERE tasks.college_id=$1
	AND tasks.category IN (%s)
	AND tasks.status='open'
	ORDER BY %s %s, id ASC
	LIMIT $2 OFFSET $3;
	`, category, filters.sortColumn(), filters.sortDirection())
	}

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

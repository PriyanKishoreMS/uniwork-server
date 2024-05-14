package data

import (
	"database/sql"
	"fmt"

	"github.com/lib/pq"
)

type TaskModel struct {
	DB *sql.DB
}

type TaskRequestModel struct {
	DB *sql.DB
}

func (t TaskModel) Create(task *Task) error {
	query := `
	INSERT INTO tasks (user_id, college_id, title, description, category, price, status, expiry, images, files)
	VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
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
		pq.Array(task.Files),
	}

	ctx, cancel := handlectx()
	defer cancel()

	return t.DB.QueryRowContext(ctx, query, args...).Scan(&task.ID)

}

func (t TaskModel) Get(id int64) (*GetTaskResponse, error) {
	query := `SELECT
  tasks.id,
  tasks.user_id,
  tasks.title,
  tasks.description,
  tasks.category,
  tasks.price,
  tasks.status,
  tasks.created_at,
  tasks.expiry,
  tasks.images,
  tasks.files,
  users.name AS user_name,
  users.avatar,
  users.rating,
  colleges.name AS college_name
FROM tasks
INNER JOIN users ON users.id = tasks.user_id
INNER JOIN colleges ON colleges.id = tasks.college_id
WHERE tasks.id = $1;
	`
	ctx, cancel := handlectx()
	defer cancel()

	var task GetTaskResponse

	dest := []interface{}{
		&task.ID,
		&task.UserID,
		&task.Title,
		&task.Description,
		&task.Category,
		&task.Price,
		&task.Status,
		&task.CreatedAt,
		&task.Expiry,
		pq.Array(&task.Images),
		pq.Array(&task.Files),
		&task.UserName,
		&task.UserAvatar,
		&task.UserRating,
		&task.CollegeName,
	}

	err := t.DB.QueryRowContext(ctx, query, id).Scan(dest...)
	if err != nil {
		return nil, err
	}

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
	fmt.Println(category, "this is the category")
	if category == "" {
		query = fmt.Sprintf(`
		SELECT
		COUNT(*) OVER () AS total,
		tasks.id, tasks.user_id, tasks.college_id, tasks.title, tasks.description, tasks.category, tasks.price, tasks.status, tasks.created_at, tasks.expiry, tasks.images, users.name, users.avatar, users.rating
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
	tasks.id, tasks.user_id, tasks.college_id, tasks.title, tasks.description, tasks.category, tasks.price, tasks.status, tasks.created_at, tasks.expiry, tasks.images, users.name, users.avatar, users.rating
	FROM tasks
	INNER JOIN users ON users.id=tasks.user_id	
	WHERE tasks.college_id=$1
	AND tasks.category='%s'
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

	for rows.Next() {
		var task TaskWithUser

		err := rows.Scan(
			&totalRecords,
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
			pq.Array(&task.Images),
			&task.UserName,
			&task.UserAvatar,
			&task.UserRating,
		)
		if err != nil {
			return nil, Metadata{}, err
		}

		tasks = append(tasks, &task)
	}

	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}
	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)

	return tasks, metadata, nil
}

func (t TaskModel) GetAllTasksOfUser(id int64, userType string, filters Filters) ([]*TaskWithUser, Metadata, error) {
	var query string
	if userType == "user" {
		fmt.Println("user")
		query = fmt.Sprintf(`
	SELECT 
	COUNT(*) OVER () AS total,
	tasks.id, tasks.user_id, tasks.college_id, tasks.title, tasks.description, tasks.category, tasks.price, tasks.status, tasks.created_at, tasks.expiry, tasks.images, users.name, users.avatar, users.rating
	FROM tasks
	INNER JOIN users ON 
	users.id=tasks.user_id	
	WHERE tasks.user_id=$1
	ORDER BY %s %s, id ASC
	LIMIT $2 OFFSET $3;
	`, filters.sortColumn(), filters.sortDirection())
	} else if userType == "worker" {
		fmt.Println("worker")
		query = fmt.Sprintf(`
	SELECT 
	COUNT(*) OVER () AS total,
	tasks.id, tasks.user_id, tasks.college_id, tasks.title, tasks.description, tasks.category, tasks.price, tasks.status, tasks.created_at, tasks.expiry, tasks.images, users.name, users.avatar, users.rating
	FROM tasks
	INNER JOIN users ON 
	users.id=tasks.worker_id	
	WHERE tasks.user_id=$1
	ORDER BY %s %s, id ASC
	LIMIT $2 OFFSET $3;
	`, filters.sortColumn(), filters.sortDirection())
	}

	ctx, cancel := handlectx()
	defer cancel()

	args := []interface{}{
		id,
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

	for rows.Next() {
		var task TaskWithUser

		err := rows.Scan(
			&totalRecords,
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
			pq.Array(&task.Images),
			&task.UserName,
			&task.UserAvatar,
			&task.UserRating,
		)
		if err != nil {
			return nil, Metadata{}, err
		}

		tasks = append(tasks, &task)
	}

	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}
	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)

	return tasks, metadata, nil
}

func (t TaskRequestModel) checkTaskRequest(userId, taskId int) (bool, error) {
	query := `SELECT EXISTS(SELECT id FROM task_requests WHERE user_id=$1 AND task_id=$2)`

	ctx, cancel := handlectx()
	defer cancel()

	var exists bool
	err := t.DB.QueryRowContext(ctx, query, userId, taskId).Scan(&exists)
	fmt.Println(exists, "This is res")
	return exists, err
}

func (t TaskRequestModel) CreateTaskRequest(userId int, taskId int) (sql.Result, error) {
	query := `INSERT INTO task_requests (task_id, user_id)
	VALUES ($1, $2)`

	ctx, cancel := handlectx()
	defer cancel()

	exists, err := t.checkTaskRequest(userId, taskId)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, fmt.Errorf("request already exists")
	}

	return t.DB.ExecContext(ctx, query, taskId, userId)
}

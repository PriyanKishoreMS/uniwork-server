package data

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/lib/pq"
)

type TaskModel struct {
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
    colleges.name AS college_name,
    COALESCE (json_agg(
        json_build_object(
			'id', task_requests.id,
			'userid', task_requests.user_id,
            'status', task_requests.status,
            'name', requesters.name,
            'avatar', requesters.avatar
        )
    ) FILTER (WHERE task_requests.id IS NOT NULL), '[]') AS requesters
FROM tasks
INNER JOIN users ON users.id = tasks.user_id
INNER JOIN colleges ON colleges.id = tasks.college_id
LEFT JOIN task_requests ON task_requests.task_id = tasks.id
LEFT JOIN users AS requesters ON requesters.id = task_requests.user_id
WHERE tasks.id = $1
GROUP BY
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
    users.name,
    users.avatar,
    users.rating,
    colleges.name;;
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
		&task.Requesters,
	}

	err := t.DB.QueryRowContext(ctx, query, id).Scan(dest...)
	if err != nil {
		return nil, err
	}

	if task.Requesters == nil {
		task.Requesters = json.RawMessage("[]")
	}

	return &task, nil
}

func (t TaskModel) GetTaskOwner(id int64) (int64, int, error) {
	query := `SELECT user_id, version FROM tasks WHERE id=$1`

	ctx, cancel := handlectx()
	defer cancel()

	var taskId int64
	var version int
	err := t.DB.QueryRowContext(ctx, query, id).Scan(&taskId, &version)
	if err != nil {
		return 0, 0, err
	}

	return taskId, version, err
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
	users.id=tasks.user_id	
	WHERE tasks.worker_id=$1
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

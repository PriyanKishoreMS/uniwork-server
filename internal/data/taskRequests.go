package data

import (
	"database/sql"
	"fmt"
)

type TaskRequestModel struct {
	DB *sql.DB
}

func (t TaskRequestModel) ApproveTaskRequest(taskId, userId int64, taskVersion int) ([3]sql.Result, error) {
	res := [3]sql.Result{}

	ctx, cancel := handlectx()
	defer cancel()

	tx, err := t.DB.BeginTx(ctx, nil)
	if err != nil {
		return res, fmt.Errorf("could not begin transaction: %w", err)
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	fmt.Println(taskVersion, "version here")

	query := `UPDATE tasks 
	SET worker_id=$1, status='assigned', version=version+1, updated_at=NOW() 
	WHERE id=$2 AND version=$3`
	result, err := tx.ExecContext(ctx, query, userId, taskId, taskVersion)
	if err != nil {
		tx.Rollback()
		if err == sql.ErrNoRows {
			return res, ErrEditConflict
		}
		return res, fmt.Errorf("could not update task: %w", err)
	}

	res[0] = result

	query = `UPDATE task_requests 
	SET status='approved', version=version+1 
	WHERE task_id=$1 AND user_id=$2`
	result, err = tx.ExecContext(ctx, query, taskId, userId)
	if err != nil {
		tx.Rollback()
		return res, fmt.Errorf("could not update task request: %w", err)
	}

	res[1] = result

	query = `DELETE FROM task_requests 
	WHERE task_id=$1 AND status <> 'approved'`
	result, err = tx.ExecContext(ctx, query, taskId)
	if err != nil {
		tx.Rollback()
		return res, fmt.Errorf("could not delete task request: %w", err)
	}

	res[2] = result

	if err = tx.Commit(); err != nil {
		return res, fmt.Errorf("could not commit task request approve transaction: %w", err)
	}

	return res, nil
}

func (t TaskRequestModel) RejectTaskRequest(taskId, userId int64) (sql.Result, error) {
	query := `UPDATE task_requests 
	SET status='rejected', version=version+1 
	WHERE task_id=$1 AND user_id=$2 AND status='pending'`

	ctx, cancel := handlectx()
	defer cancel()

	res, err := t.DB.ExecContext(ctx, query, taskId, userId)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (t TaskRequestModel) CheckTaskRequestStatus(taskId int64) (bool, error) {
	query := `SELECT EXISTS(SELECT id FROM task_requests WHERE task_id=$1 AND status='approved')`

	ctx, cancel := handlectx()
	defer cancel()

	var exists bool
	err := t.DB.QueryRowContext(ctx, query, taskId).Scan(&exists)
	return exists, err
}

func (t TaskRequestModel) CreateTaskRequest(userId, taskId int64) (sql.Result, error) {
	query := `INSERT INTO task_requests (task_id, user_id)
	VALUES ($1, $2)
	ON CONFLICT(task_id, user_id) DO NOTHING
	`

	ctx, cancel := handlectx()
	defer cancel()

	return t.DB.ExecContext(ctx, query, taskId, userId)
}

func (t TaskRequestModel) DeleteTaskRequest(userId, taskId int64) (sql.Result, error) {
	query := `DELETE FROM task_requests WHERE user_id=$1 AND task_id=$2 AND status='pending'`

	ctx, cancel := handlectx()
	defer cancel()

	return t.DB.ExecContext(ctx, query, userId, taskId)
}

type checkoutData struct {
	TaskId        int64  `json:"task_id"`
	WorkerId      int64  `json:"worker_id" `
	Title         string `json:"title" `
	Category      string `json:"category" `
	Price         int64  `json:"price" `
	CreatedAt     string `json:"created_at" `
	Expiry        string `json:"expiry" `
	WorkerName    string `json:"worker_name" `
	WorkerCollege string `json:"worker_college" `
	WorkerAvatar  string `json:"worker_avatar" `
}

func (t TaskRequestModel) GetCheckoutTaskRequest(userId, taskId int64) (*checkoutData, error) {
	query := `SELECT t.id, t.title, t.category, t.price, t.created_at, t.expiry, u.id, u.name, c.name, u.avatar
	FROM tasks t
	JOIN users u ON u.id = $1 AND t.id=$2
	JOIN colleges c ON c.id = u.college_id;
	`

	ctx, cancel := handlectx()
	defer cancel()

	row := t.DB.QueryRowContext(ctx, query, userId, taskId)

	var data checkoutData

	err := row.Scan(&data.TaskId, &data.Title, &data.Category, &data.Price, &data.CreatedAt, &data.Expiry, &data.WorkerId, &data.WorkerName, &data.WorkerCollege, &data.WorkerAvatar)

	if err != nil {
		return nil, err
	}

	return &data, nil

}

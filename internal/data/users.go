package data

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/labstack/gommon/log"
)

type User struct {
	ID             int64     `json:"id"`
	CollegeID      int64     `json:"college_id"`
	Name           string    `json:"name" validate:"required"`
	Email          string    `json:"email,omitempty" validate:"required,email"`
	Mobile         string    `json:"mobile,omitempty"`
	Avatar         string    `json:"avatar,omitempty"`
	Dept           string    `json:"dept" validate:"required"`
	TasksCompleted int       `json:"tasks_completed,omitempty"`
	Earned         int64     `json:"earned,omitempty"`
	Rating         float64   `json:"rating,omitempty"`
	CreatedAt      time.Time `json:"created_at,omitempty"`
	Version        int       `json:"version,omitempty"`
}

type UserModel struct {
	DB *sql.DB
}

func (u UserModel) Register(user *User) (int64, error) {
	query := `
	INSERT INTO users (college_id, name, email, dept, mobile)
	VALUES(?, ?, ?, ?, ?)
	`
	args := []interface{}{
		user.CollegeID,
		user.Name,
		user.Email,
		user.Dept,
		user.Mobile,
	}

	ctx, cancel := handlectx()
	defer cancel()

	res, err := u.DB.ExecContext(ctx, query, args...)
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

func (u UserModel) Get(id int64) (*User, error) {
	query := `SELECT id, college_id, name, email, mobile, dept, avatar, tasks_completed, earned, rating, created_at, version
	FROM users
	WHERE id=?
	`
	ctx, cancel := handlectx()
	defer cancel()

	var user User

	dest := []interface{}{
		&user.ID,
		&user.CollegeID,
		&user.Name,
		&user.Email,
		&user.Mobile,
		&user.Dept,
		&user.Avatar,
		&user.TasksCompleted,
		&user.Earned,
		&user.Rating,
		&user.CreatedAt,
		&user.Version,
	}

	err := u.DB.QueryRowContext(ctx, query, id).Scan(dest...)
	if err != nil {
		return nil, ErrRecordNotFound
	}

	return &user, nil
}

func (u UserModel) GetUserByEmail(email string) (*User, error) {
	query := `SELECT id, college_id, name, email, mobile, dept, avatar, tasks_completed, earned, rating, created_at, version
	FROM users
	WHERE email=?
	`
	ctx, cancel := handlectx()
	defer cancel()

	var user User

	dest := []interface{}{
		&user.ID,
		&user.CollegeID,
		&user.Name,
		&user.Email,
		&user.Mobile,
		&user.Dept,
		&user.Avatar,
		&user.TasksCompleted,
		&user.Earned,
		&user.Rating,
		&user.CreatedAt,
		&user.Version,
	}

	err := u.DB.QueryRowContext(ctx, query, email).Scan(dest...)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (u UserModel) Update(user *User) (int64, error) {
	query := `UPDATE users 
	SET college_id=?, name=?, email=?, mobile=?, dept=?, avatar=?, version=version+1
	WHERE id=? AND version=?
	`

	args := []interface{}{
		&user.CollegeID,
		&user.Name,
		&user.Email,
		&user.Mobile,
		&user.Dept,
		&user.Avatar,
		&user.ID,
		&user.Version,
	}

	ctx, cancel := handlectx()
	defer cancel()

	res, err := u.DB.ExecContext(ctx, query, args...)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return 0, ErrEditConflict
		default:
			return 0, err
		}
	}

	RowsAffected, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}
	return RowsAffected, nil
}

func (u UserModel) Delete(id int64) (int64, error) {
	query := `
	DELETE FROM users WHERE id=?
	`
	ctx, cancel := handlectx()
	defer cancel()

	res, err := u.DB.ExecContext(ctx, query, id)
	if err != nil {
		return 0, err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}

	if rowsAffected == 0 {
		return 0, ErrRecordNotFound
	}

	return rowsAffected, nil
}

func (u UserModel) GetAllInCollege(name string, college_id int64, filters Filters) ([]*User, Metadata, error) {
	query := fmt.Sprint(`
	SELECT COUNT(*) OVER () AS total,
	id, college_id, name, dept, avatar, rating
	FROM users
	WHERE LOWER(name) LIKE LOWER(CONCAT('%', ? ,'%'))
	AND college_id=?
	ORDER BY ` + filters.sortColumn() + " " + filters.sortDirection() + `, id ASC
	LIMIT ? OFFSET ?;
	`)

	ctx, cancel := handlectx()
	defer cancel()

	args := []interface{}{
		name,
		college_id,
		filters.limit(),
		filters.offset(),
	}

	rows, err := u.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, Metadata{}, err
	}

	defer rows.Close()

	totalRecords := 0
	users := []*User{}

	for rows.Next() {
		var user User

		err := rows.Scan(
			&totalRecords,
			&user.ID,
			&user.CollegeID,
			&user.Name,
			&user.Dept,
			&user.Avatar,
			&user.Rating,
		)
		if err != nil {
			return nil, Metadata{}, err
		}

		users = append(users, &user)
	}

	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}
	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)

	return users, metadata, nil
}

package data

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

type College struct {
	ID      int64  `json:"id"`
	Name    string `json:"name" validate:"required"`
	Domain  string `json:"domain" validate:"required,email"`
	Version string `json:"version"`
}

type CollegeModel struct {
	DB *sql.DB
}

func handlectx() (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	return ctx, cancel
}

func (c CollegeModel) Create(college *College) (int64, error) {
	query := `INSERT INTO colleges (name, domain)
	VALUES(?, ?)
	`

	args := []interface{}{
		college.Name,
		college.Domain,
	}

	ctx, cancel := handlectx()
	defer cancel()

	res, err := c.DB.ExecContext(ctx, query, args...)
	if err != nil {
		return 0, err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}

	if rowsAffected == 0 {
		return 0, sql.ErrNoRows
	}

	LastInsertId, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return LastInsertId, nil
}

func (c CollegeModel) Get(id int64) (*College, error) {
	query := `SELECT id, name, domain, version
	FROM colleges
	WHERE id=?
	`
	ctx, cancel := handlectx()
	defer cancel()

	var college College

	err := c.DB.QueryRowContext(ctx, query, id).Scan(&college.ID, &college.Name, &college.Domain, &college.Version)
	if err != nil {
		return nil, err
	}

	return &college, nil
}

func (c CollegeModel) Update(college *College) (int64, error) {

	query := `UPDATE colleges 
	SET name=?, domain=?, version=version+1
	WHERE id=? AND version=?
	`

	args := []interface{}{
		college.Name,
		college.Domain,
		college.ID,
		college.Version,
	}

	ctx, cancel := handlectx()
	defer cancel()

	res, err := c.DB.ExecContext(ctx, query, args...)
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

func (c CollegeModel) Delete(id int64) (int64, error) {
	query := `
	DELETE FROM colleges WHERE id=?
	`
	ctx, cancel := handlectx()
	defer cancel()

	res, err := c.DB.ExecContext(ctx, query, id)
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

package data

import (
	"database/sql"
	"errors"
	"fmt"
)

type College struct {
	ID      int64  `json:"id"`
	Name    string `json:"name" validate:"required"`
	Domain  string `json:"domain" validate:"required,email"`
	Version int    `json:"version"`
}

type CollegeModel struct {
	DB *sql.DB
}

func (c CollegeModel) Create(college *College) error {
	query := `INSERT INTO colleges (name, domain)
	VALUES($1, $2)
	RETURNING id
	`

	args := []interface{}{
		college.Name,
		college.Domain,
	}

	ctx, cancel := handlectx()
	defer cancel()

	return c.DB.QueryRowContext(ctx, query, args...).Scan(&college.ID)

}

func (c CollegeModel) Get(id int64) (*College, error) {
	query := `SELECT id, name, domain, version
	FROM colleges
	WHERE id=$1
	`
	ctx, cancel := handlectx()
	defer cancel()

	var college College

	dest := []interface{}{
		&college.ID,
		&college.Name,
		&college.Domain,
		&college.Version,
	}

	err := c.DB.QueryRowContext(ctx, query, id).Scan(dest...)
	if err != nil {
		return nil, err
	}

	return &college, nil
}

func (c CollegeModel) Update(college *College) error {

	query := `UPDATE colleges 
	SET name=$1, domain=$2, version=version+1
	WHERE id=$3 AND version=$4
	RETURNING id
	`

	args := []interface{}{
		college.Name,
		college.Domain,
		college.ID,
		college.Version,
	}

	ctx, cancel := handlectx()
	defer cancel()

	err := c.DB.QueryRowContext(ctx, query, args...).Scan(&college.ID)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		default:
			return err
		}
	}

	return nil

}

func (c CollegeModel) Delete(id int64) error {
	query := `
	DELETE FROM colleges WHERE id=$1
	RETURNING id
	`
	ctx, cancel := handlectx()
	defer cancel()

	err := c.DB.QueryRowContext(ctx, query, id).Scan(&id)
	if err != nil {
		return err
	}

	return nil
}

func (c CollegeModel) GetAll(name string, filters Filters) ([]*College, Metadata, error) {
	query := fmt.Sprint(`
	SELECT COUNT(*) OVER () AS total,
	id, name, domain, version
	FROM colleges
	WHERE name ILIKE '%' || $1 || '%' 
	ORDER BY ` + filters.sortColumn() + " " + filters.sortDirection() + `, id ASC
	LIMIT $2 OFFSET $3;
	`)

	ctx, cancel := handlectx()
	defer cancel()

	args := []interface{}{
		name,
		filters.limit(),
		filters.offset(),
	}

	rows, err := c.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, Metadata{}, err
	}

	defer rows.Close()

	totalRecords := 0
	colleges := []*College{}

	for rows.Next() {
		var college College

		err := rows.Scan(
			&totalRecords,
			&college.ID,
			&college.Name,
			&college.Domain,
			&college.Version,
		)
		if err != nil {
			return nil, Metadata{}, err
		}

		colleges = append(colleges, &college)
	}

	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}
	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)

	return colleges, metadata, nil

}

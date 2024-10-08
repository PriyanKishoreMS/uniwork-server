package data

import (
	"database/sql"
	"errors"
	"fmt"
)

type UserModel struct {
	DB *sql.DB
}

func (u UserModel) Register(user *User) error {
	query := `
	INSERT INTO users (college_id, name, email, dept, mobile)
	VALUES($1, $2, $3, $4, $5)
	RETURNING id, name, email, dept, mobile
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

	return u.DB.QueryRowContext(ctx, query, args...).Scan(&user.ID, &user.Name, &user.Email, &user.Dept, &user.Mobile)
}

func (u UserModel) Get(id int64) (*User, error) {
	query := `select users.id, users.college_id, users.name as user_name, users.email, users.avatar, users.dept, users.tasks_completed, users.earned, users.rating, colleges.name as college_name from users inner join colleges on colleges.id = users.college_id where users.id=$1;
	`
	ctx, cancel := handlectx()
	defer cancel()

	var user User

	dest := []interface{}{
		&user.ID,
		&user.CollegeID,
		&user.Name,
		&user.Email,
		&user.Avatar,
		&user.Dept,
		&user.TasksCompleted,
		&user.Earned,
		&user.Rating,
		&user.CollegeName,
	}

	err := u.DB.QueryRowContext(ctx, query, id).Scan(dest...)
	if err != nil {
		return nil, ErrRecordNotFound
	}

	return &user, nil
}

func (u UserModel) GetUserByEmail(email string) (*User, error) {
	query := `select users.id, users.college_id, users.name as user_name, users.email, users.avatar, users.dept, users.tasks_completed, users.earned, users.rating, colleges.name as college_name from users inner join colleges on colleges.id = users.college_id where users.email=$1
	`
	ctx, cancel := handlectx()
	defer cancel()

	var user User

	dest := []interface{}{
		&user.ID,
		&user.CollegeID,
		&user.Name,
		&user.Email,
		&user.Avatar,
		&user.Dept,
		&user.TasksCompleted,
		&user.Earned,
		&user.Rating,
		&user.CollegeName,
	}

	err := u.DB.QueryRowContext(ctx, query, email).Scan(dest...)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (u UserModel) Update(user *User) error {
	query := `UPDATE users 
	SET college_id=$1, name=$2, email=$3, mobile=$4, dept=$5, avatar=$6, version=version+1
	WHERE id=$7 AND version=$8
	RETURNING id
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

	err := u.DB.QueryRowContext(ctx, query, args...).Scan(&user.ID)
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

func (u UserModel) Delete(id int64) error {
	query := `
	DELETE FROM users WHERE id=$1
	RETURNING id
	`
	ctx, cancel := handlectx()
	defer cancel()

	err := u.DB.QueryRowContext(ctx, query, id).Scan(&id)
	if err != nil {
		return err
	}

	return nil

}

func (u UserModel) GetAllInCollege(name string, college_id int64, filters Filters) ([]*User, Metadata, error) {
	query := fmt.Sprint(`
	SELECT COUNT(*) OVER () AS total,
	id, college_id, name, dept, avatar, rating
	FROM users
	WHERE name ILIKE '%' || $1 || '%'
	AND college_id=$2
	ORDER BY ` + filters.sortColumn() + " " + filters.sortDirection() + `, id ASC
	LIMIT $3 OFFSET $4;
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

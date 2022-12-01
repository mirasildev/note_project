package postgres

import (
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/mirasildev/note_project/storage/repo"
)

type userRepo struct {
	db *sqlx.DB
}

func NewUser(db *sqlx.DB) repo.UserStorageI {
	return &userRepo{
		db: db,
	}
}

func (ur *userRepo) Create(u *repo.User) (*repo.User, error) {
	query := `
		INSERT INTO users(
			first_name,
			last_name,
			phone_number,
			email,
			image_url,
			created_at	                  
		) VALUES($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at 
	`

	err := ur.db.QueryRow(
		query,
		u.FirstName,
		u.LastName,
		u.PhoneNumber,
		u.Email,
		u.ImageURL,
		u.CreatedAt,
	).Scan(
		&u.ID,
		&u.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return u, nil
}

func (ur *userRepo) Get(id int64) (*repo.User, error) {
	var result repo.User

	query := `
		SELECT 
			id, 
			last_name,
			phone_number,
			email,
			image_url,
			created_at
		FROM users
		WHERE id=$1
	`

	row := ur.db.QueryRow(query,id)
	err := row.Scan(
		&result.ID,
		&result.FirstName,
		&result.LastName,
		&result.Email,
		&result.ImageURL,
		&result.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (ur *userRepo) GetAllUsers(params *repo.GetAllUsersParams) (*repo.GetAllUsersResult, error) {
	result := repo.GetAllUsersResult{
		Users: make([]*repo.User, 0),
	}

	offset := (params.Page - 1) * params.Limit

	limit := fmt.Sprintf(" LIMIT %d OFFSET %d ", params.Limit, offset)

	filter := "" 
	if params.Search != "" {
		str := "%" + params.Search + "%"
		filter += fmt.Sprintf(
			`WHERE first_name ILIKE '%s' OR last_name ILIKE '%s' 
			OR email ILIKE '%s' OR phone_number ILIKE '%s' `,
			str, str, str, str,
		) 
	}

	query := `
		SELECT 
			id, 
			first_name,
			last_name,
			phone_number,
			email,
			image_url,
			created_at
		FROM users
		` + filter + `
		ORDER BY created_at desc
		` + limit
	fmt.Println(query)

	rows, err := ur.db.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var u repo.User

		err := rows.Scan(
			&u.ID,
			&u.FirstName,
			&u.LastName,
			&u.PhoneNumber,
			&u.Email,
			&u.ImageURL,
			&u.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		result.Users = append(result.Users, &u)
	}

	queryCount := `SELECT count(1) FROM users` + filter
	err = ur.db.QueryRow(queryCount).Scan(&result.Count)
	if err != nil {
		return nil, err
	}
	fmt.Println(result.Users)
	return &result, nil
}

func (ur *userRepo) Update(u *repo.User) (*repo.User, error) {
	query := `
		UPDATE users SET 
			first_name=$1,
			last_name=$2,
			phone_number=$3,
			email=$4,
			image_url=$5
		WHERE id=$6
		RETURNING id, first_name, last_name, phone_number, email,
		image_url, created_at
	`
	var result repo.User
	err:= ur.db.QueryRow(
		query,
		u.FirstName,
		u.LastName,
		u.PhoneNumber,
		u.Email,
		u.ImageURL,
		u.ID,
	).Scan(
		&result.ID,
		&result.FirstName,
		&result.LastName,
		&result.PhoneNumber,
		&result.Email,
		&result.ImageURL,
		&result.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (ur *userRepo) Delete(id int64) error {
	queryDeleteNotes := "DELETE FROM notes WHERE user_id=$1"
	_, err := ur.db.Exec(queryDeleteNotes, id)
	if err != nil {
		return err
	}

	query := "DELETE FROM users WHERE id=$1"
	result, err := ur.db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsCount, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsCount == 0 {
		return sql.ErrNoRows
	}

	return nil
}
 
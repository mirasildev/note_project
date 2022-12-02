package postgres

import (
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/mirasildev/note_project/storage/repo"
)

type noteRepo struct {
	db *sqlx.DB
}

func NewNote(db *sqlx.DB) repo.NoteStorageI {
	return &noteRepo{
		db: db,
	}
}

func (nt *noteRepo) Create(n *repo.Note) (*repo.Note, error) {
	query := `
		INSERT INTO notes(
			user_id,
			title,
			description,
			created_at,
			updated_at,
			deleted_at
		) VALUES($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at
	`

	err := nt.db.QueryRow(
		query,
		n.UserID,
		n.Title,
		n.Description,
		n.CreatedAt,
		n.UpdatedAt,
		n.DeletedAt,
	).Scan(&n.ID, &n.CreatedAt)
	if err != nil {
		return nil, err
	}

	return n, nil
}

func (nt *noteRepo) Get(id int64) (*repo.Note, error) {
	var result repo.Note

	query := `
		SELECT 
			id,
			user_id,
			title,
			description,
			created_at,
			updated_at
		FROM notes
		WHERE id=$1 AND deleted_at IS NULL
	`

	row := nt.db.QueryRow(query, id)
	err := row.Scan(
		&result.ID,
		&result.UserID,
		&result.Title,
		&result.Description,
		&result.CreatedAt,
		&result.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (nt *noteRepo) GetAllNotes(params *repo.GetAllNotesParams) (*repo.GetAllNotesResult, error) {
	result := repo.GetAllNotesResult{
		Notes: make([]*repo.Note, 0),
	}

	offset := (params.Page - 1) * params.Limit
	limit := fmt.Sprintf(" LIMIT %d OFFSET %d", params.Limit, offset)
	filter := ""
	if params.Search != "" {
		str := "%" + fmt.Sprint(params.UserID) + "%"
		filter += fmt.Sprintf(
			` WHERE user_id='%s' AND deleted_at IS NULL`,
			str,
		)
	}

	query := `
		SELECT 
			id,
			user_id,
			title,
			description,
			created_at,
			updated_at
		FROM notes
		` + filter + `
		ORDER BY created_at desc
		` + limit

	rows, err := nt.db.Query(query)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var note repo.Note

		err := rows.Scan(
			&note.ID,
			&note.UserID,
			&note.Title,
			&note.Description,
			&note.CreatedAt,
			&note.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		result.Notes = append(result.Notes, &note)
	}

	return &result, nil
}

func (nt *noteRepo) Update(n *repo.Note) (*repo.Note, error) {
	query := `
		UPDATE notes SET
			user_id=$1,
			title=$2,
			description=$3,
			updated_at=$4
		WHERE id=$5 AND deleted_at IS NULL
		RETURNING id, user_id, title, description, created_at, updated_at
	`

	var result repo.Note

	err := nt.db.QueryRow(query,
		n.UserID,
		n.Title,
		n.Description,
		n.UpdatedAt,
		n.ID,
	).Scan(
		&result.ID,
		&result.UserID,
		&result.Title,
		&result.Description,
		&result.CreatedAt,
		&result.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (nt *noteRepo) Delete(id int64) error {

	query := "DELETE FROM notes WHERE id=$1"
	result, err := nt.db.Exec(query, id)
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

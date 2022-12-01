package repo

import "time"

type Note struct {
	ID          int64
	UserID      int64
	Title       string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time
}

type GetAllNotesParams struct {
	UserID string
	Limit int32
	Page int32
	Search string
}

type GetAllNotesResult struct {
	Notes []*Note
	Count int32
}

type NoteStorageI interface {
	Create(n *Note) (*Note, error)
	Get(id int64) (*Note, error)
	GetAllNotes(params *GetAllNotesParams) (*GetAllNotesResult, error)
	Update(n *Note) (*Note, error)
	Delete(id int64) error
}

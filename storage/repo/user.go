package repo

import "time"

type User struct {
	ID          int64
	FirstName   string
	LastName    string
	PhoneNumber *string
	Email       string
	Password    string
	ImageURL    *string
	CreatedAt   time.Time
}

type GetAllUsersParams struct {
	Limit  int32
	Page   int32
	Search string
}

type GetAllUsersResult struct {
	Count int32
	Users []*User
}

type UserStorageI interface {
	Create(n *User) (*User, error)
	Get(id int64) (*User, error)
	GetAllUsers(params *GetAllUsersParams) (*GetAllUsersResult, error)
	Update(n *User) (*User, error)
	Delete(id int64) error
	GetByEmail(email string) (*User, error)
}

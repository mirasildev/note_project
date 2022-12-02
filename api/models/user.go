package models

import "time"

type User struct {
	ID          int64     `json:"id"`
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	PhoneNumber *string   `json:"phone_number"`
	Email       string    `json:"email"`
	ImageURL    *string   `json:"image_url"`
	CreatedAt   time.Time `json:"created_at"`
}

type CreateUserRequest struct {
	FirstName   string  `json:"first_name" binding:"required,min=2,max=30"`
	LastName    string  `json:"last_name" binding:"required,min=2,max=30"`
	PhoneNumber *string `json:"phone_number"` // *
	Email       string  `json:"email" binding:"required,email"`
	ImageURL    *string `json:"image_url"` // *
}

type GetAllUsersResponse struct {
	Users []*User `json:"users"`
	Count int32   `json:"count"`
}

type VerifyRequest struct {
	Email string `json:"email" binding:"required,email"`
	Code  string `json:"code" binding:"required"`
}

type UpdateUserRequest struct {
	// ID          int64   `json:"id"`
	FirstName   string    `json:"first_name" binding:"required,min=2,max=30"`
	LastName    string    `json:"last_name" binding:"required,min=2,max=30"`
	PhoneNumber *string   `json:"phone_number"` // *
	Email       string    `json:"email" binding:"required,email"`
	ImageURL    *string   `json:"image_url"` // *
	// UpdatedAt   time.Time `json:"updated_at"`
}

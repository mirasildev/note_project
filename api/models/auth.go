package models

import "time"

type RegisterRequest struct {
	FirstName string `json:"first_name" binding:"required,min=2,max=30"`
	LastName  string `json:"last_name" binding:"required,min=2,max=30"`
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=6,max=16"`
}

type AuthResponse struct {
	ID          int64     `json:"id"`
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	Email       string    `json:"email"`
	Username    string    `json:"username"`
	CreatedAt   time.Time `json:"created_at"`
	AccessToken string    `json:"access_token"`
}

type LoginRequest struct {
	Email string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6,max=16"`
}

type UpdatePassword struct {
	Password string `json:"password" binding:"required"`
}
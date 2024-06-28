package model

import "time"

type User struct {
	ID           string    `json:"id"`
	Login        string    `json:"login" binding:"required,min=3,max=16"`
	Password     string    `json:"password" binding:"required,min=8,max=32"`
	Role         string    `json:"role"`
	DateAdded    time.Time `json:"dateAdded"`
}

func (u *User) DTO() *User {
	return &User{
		ID: u.ID,
		Login: u.Login,
		// IGNORING u.Password
		Role: u.Role,
		DateAdded: u.DateAdded,
	}
}

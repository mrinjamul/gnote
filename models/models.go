package models

import "time"

// Note struct
type Note struct {
	ID        uint64    `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	User      string    `json:"user"`
	Archived  bool      `json:"archived"`
	CreatedAt time.Time `json:"createdat"`
	UpdatedAt time.Time `json:"updatedat"`
}

// User struct
type User struct {
	ID        uint64    `json:"id" gorm:"primary_key"`
	UserID    string    `json:"user_id" gorm:"default:uuid_generate_v3()"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"createdat"`
	UpdatedAt time.Time `json:"updatedat"`
}

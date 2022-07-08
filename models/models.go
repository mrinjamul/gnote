package models

import (
	"database/sql"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

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

// User is a user of the application
type User struct {
	ID         uint         `json:"id" gorm:"primary_key,autoIncrement,not null"`
	FirstName  string       `json:"first_name" gorm:"not null"`
	MiddleName string       `json:"middle_name,omitempty"`
	LastName   string       `json:"last_name" gorm:"not null"`
	UserName   string       `json:"username"  gorm:"unique,not null"`
	Email      string       `json:"email" gorm:"unique"`
	DOB        time.Time    `json:"dob" gorm:"not null"`
	Password   string       `json:"password" gorm:"not null"`
	Role       string       `json:"role" gorm:"not null"`
	Level      int          `json:"level" gorm:"not null"`
	CreatedAt  time.Time    `json:"created_at" gorm:"not null"`
	UpdatedAt  time.Time    `json:"updated_at" gorm:"not null"`
	DeletedAt  sql.NullTime `json:"deleted_at" gorm:"index,not null"`
}

// Create a struct that models the structure of a user in the request body
type Credentials struct {
	UserName string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Claims
type Claims struct {
	UserName string `json:"username"`
	Role     string `json:"role"`
	Level    int    `json:"level"`
	jwt.RegisteredClaims
}

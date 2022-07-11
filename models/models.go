package models

import (
	"database/sql"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// Note struct
type Note struct {
	ID        uint64    `json:"id" gorm:"primary_key,autoIncrement,not null"`
	Title     string    `json:"title,omitempty"`
	Content   string    `json:"content" gorm:"not null"`
	Username  string    `json:"username" gorm:"not null"`
	Archived  bool      `json:"archived,omitempty"`
	CreatedAt time.Time `json:"created_at" gorm:"not null"`
	UpdatedAt time.Time `json:"updated_at" gorm:"index,not null"`
}

// User is a user of the application
type User struct {
	ID         uint         `json:"id" gorm:"primary_key,autoIncrement,not null"`
	FirstName  string       `json:"first_name" gorm:"not null"`
	MiddleName string       `json:"middle_name,omitempty"`
	LastName   string       `json:"last_name" gorm:"not null"`
	Username   string       `json:"username"  gorm:"unique,not null"`
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
	Username string `json:"username,omitempty"`
	Email    string `json:"email,omitempty"`
	Password string `json:"password"`
}

// Claims
type Claims struct {
	Username string `json:"username"`
	Role     string `json:"role"`
	Level    int    `json:"level"`
	jwt.RegisteredClaims
}

// Config is the configuration for CLI
type Config struct {
	Username string `json:"username,omitempty"`
	Token    string `json:"token,omitempty"`
	APIToken string `json:"api_token,omitempty"`
}

package utils

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
	"unicode"

	"github.com/mrinjamul/gnote/models"
	"golang.org/x/crypto/bcrypt"
)

var ApiURL = "https://gnote.up.railway.app"

// sendRequest sends a request to the API
func sendRequest(method, path string, jsonData []byte) (string, error) {
	// Create a new request
	req, err := http.NewRequest(method, ApiURL+path, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}
	// Set the request's header
	req.Header.Set("Content-Type", "application/json")
	// Set the request's parameters
	// q := req.URL.Query()
	// for k, v := range params {
	// 	q.Add(k, v)
	// }
	// req.URL.RawQuery = q.Encode()
	// Set Bearer authorization
	req.Header.Set("Authorization", "Bearer "+"token")
	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	// Read the response
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

// CreateNote creates a note
func CreateNote(title, content string) (string, error) {
	jsonStr := []byte(`{"title":"` + title + `", "content":"` + content + `"}`)
	return sendRequest("POST", "/api/notes", jsonStr)
}

// GetNotes gets all notes
func GetNotes() (string, error) {
	return sendRequest("GET", "/api/notes", nil)
}

// GetNote gets a note
func GetNote(id string) (string, error) {
	jsonStr := []byte(`{"id":"` + id + `"}`)
	return sendRequest("GET", "/api/notes/"+id, jsonStr)
}

// UpdateNote updates a note
func UpdateNote(id, title, content string) (string, error) {
	jsonStr := []byte(`{"title":"` + title + `", "content":"` + content + `"}`)
	return sendRequest("PUT", "/api/notes/"+id, jsonStr)
}

// DeleteNote deletes a note
func DeleteNote(id string) (string, error) {
	jsonStr := []byte(`{"id":"` + id + `"}`)
	return sendRequest("DELETE", "/api/notes/"+id, jsonStr)
}

// GetEnv gets the environment variable
func GetEnv(key string) string {
	return os.Getenv(key)
}

// ParseToken parses the token from authorization header
func ParseToken(authorization string) (string, error) {
	if strings.HasPrefix(authorization, "Bearer ") {
		return strings.TrimPrefix(authorization, "Bearer "), nil
	}
	return "", fmt.Errorf("invalid authorization header")
}

// ToMaxAge converts the duration to max age
func ToMaxAge(expire time.Time) int {
	maxAge := time.Until(expire).Seconds()
	return int(maxAge)
}

// VerifyHash verifies the hashed password
func VerifyHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		return err == nil
	}
	return true
}

// ValidateUser the user informations
func ValidateUser(u *models.User) error {
	if u.FirstName == "" || u.LastName == "" {
		return errors.New("first name and last name are required")
	}
	if u.UserName == "" {
		return errors.New("username is required")
	}
	if u.Email == "" {
		return errors.New("email is required")
	}
	if u.Password == "" {
		return errors.New("password is required")
	}

	// check if password contains at least 7 characters,
	// at least one letter, one number and one special character
	if !IsValidPassword(u.Password) {
		return errors.New("password should contain at least seven characters, one number and one special character")
	}
	return nil
}

// IsRestrictedUser checks if the user is restricted
func IsRestrictedUser(username string) bool {
	restrictedList := []string{"admin", "root", "me", "system", "search"}
	for _, restrictedUser := range restrictedList {
		if username == restrictedUser {
			return true
		}
	}
	return false
}

// IsValidUserName checks if the username is valid
func IsValidUserName(username string) bool {
	if username == "" {
		return false
	}
	// username can be only alphanumeric
	for _, char := range username {
		if !unicode.IsLetter(char) && !unicode.IsNumber(char) {
			return false
		}
	}

	if IsRestrictedUser(username) {
		return false
	}

	// return true
	return len(username) >= 3
}

// IsValidPassword checks if the password (Strength) is valid
func IsValidPassword(s string) bool {
	var (
		hasMinLen  = false
		hasUpper   = false
		hasLower   = false
		hasNumber  = false
		hasSpecial = false
	)
	if len(s) >= 7 {
		hasMinLen = true
	}
	for _, char := range s {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}
	return hasMinLen && hasUpper && hasLower && hasNumber && hasSpecial
}

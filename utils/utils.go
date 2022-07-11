package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/mrinjamul/gnote/models"
	"golang.org/x/crypto/bcrypt"
)

// var ApiURL = "https://gnote.up.railway.app"
// var ApiURL = "https://note.mrinjamul.in"
var ApiURL = "http://localhost:8080"

// HomeDir returns the home directory of the current user
func HomeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	home, _ := os.UserHomeDir()
	return home
}

func GetConfig() (*models.Config, error) {
	// Get Home location
	home := HomeDir()
	// Get config file path
	// using filepath join
	configFilePath := filepath.Join(home, ".gnote")
	configFile := filepath.Join(configFilePath, "config.json")
	// check if the file exists
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		// Create the config file
		createIfNotExist(configFile, configFilePath)
	}

	var config models.Config
	// Read the config file and unmarshal it to config
	configFileContent, err := ioutil.ReadFile(configFile)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(configFileContent, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}

// SaveToken saves the token to the config file
func SaveToken(token string) error {
	// Get Home location
	home := HomeDir()
	// Get config file path
	// using filepath join
	configFilePath := filepath.Join(home, ".gnote")
	configFile := filepath.Join(configFilePath, "config.json")
	// check if the file exists
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		// Create the config file
		createIfNotExist(configFile, configFilePath)
	}
	// Read the config file and unmarshal it to config
	configFileContent, err := ioutil.ReadFile(configFile)
	if err != nil {
		panic(err)
	}
	var config models.Config
	err = json.Unmarshal(configFileContent, &config)
	if err != nil {
		panic(err)
	}
	config.Token = token
	// Marshal the config to json
	jsonStr, err := json.Marshal(config)
	if err != nil {
		panic(err)
	}
	// Write the json to the config file
	err = ioutil.WriteFile(configFile, jsonStr, os.ModePerm)
	if err != nil {
		panic(err)
	}
	return nil
}

// createIfNotExist creates a file if it doesn't exist
func createIfNotExist(file string, path string) {
	// Check if directory exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		// Create the directory
		os.MkdirAll(path, os.ModePerm)
	}
	// Create the config file
	os.Create(file)
}

// GenTips generates random tips
func GenTips() string {
	tips := []string{
		"Use `gnote login` to login to gnote",
		"Use `gnote logout` to logout from gnote",
		"Use `gnote create` to create a new note",
		"Use `gnote read [id]` to read a note",
		"Use `gnote update [note]` to update a note",
		"Use `gnote delete [note]` to delete a note",
		"Use `gnote list` to list all notes",
		"Use `gnote search [query]` to search notes",
		"Use `gnote version` to get gnote's version",
		"Use `gnote serve` to start gnote backend server",
		"Use `gnote help` to show this message",
		"Use `gnote help [command]` to show help for a command",
	}
	rand.Seed(time.Now().UnixNano())
	return tips[rand.Intn(len(tips))]
}

// CLILogin logs in to the API
func CLILogin(username, password string) ([]byte, error) {
	jsonStr := []byte(`{"username":"` + username + `", "password":"` + password + `"}`)
	body, err := sendRequest("POST", "/auth/login", jsonStr, "")
	if err != nil {
		return nil, err
	}
	return body, nil
}

// CLISignup signs up to the API
func CLISignup(username, password string) ([]byte, error) {
	jsonStr := []byte(`{"username":"` + username + `", "password":"` + password + `"}`)
	body, err := sendRequest("POST", "/auth/signup", jsonStr, "")
	if err != nil {
		return nil, err
	}
	return body, nil
}

// sendRequest sends a request to the API
func sendRequest(method, path string, jsonData []byte, token string) ([]byte, error) {
	// Create a new request
	req, err := http.NewRequest(method, ApiURL+path, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
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
	req.Header.Set("Authorization", "Bearer "+token)
	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	// Read the response
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

// CreateNote creates a note
func CreateNote(title, content string, token string) (models.Note, error) {
	var note models.Note
	jsonStr := []byte(`{"title":"` + title + `", "content":"` + content + `"}`)
	body, err := sendRequest("POST", "/api/notes", jsonStr, token)
	if err != nil {
		return note, err
	}
	// Parse json body
	err = json.Unmarshal(body, &note)
	if err != nil {
		return note, err
	}
	return note, nil
}

// GetNotes gets all notes
func GetNotes(token string) ([]models.Note, error) {
	var notes []models.Note
	body, err := sendRequest("GET", "/api/notes", nil, token)
	if err != nil {
		return notes, err
	}
	// Parse json body
	err = json.Unmarshal(body, &notes)
	if err != nil {
		return notes, err
	}
	return notes, err
}

// GetNote gets a note
func GetNote(id string, token string) (models.Note, error) {
	var note models.Note
	jsonStr := []byte(`{"id":"` + id + `"}`)
	body, err := sendRequest("GET", "/api/notes/"+id, jsonStr, token)
	if err != nil {
		return note, err
	}
	// Parse json body
	err = json.Unmarshal(body, &note)
	if err != nil {
		return note, err
	}
	return note, nil
}

// UpdateNote updates a note
func UpdateNote(id, title, content string, token string) (models.Note, error) {
	var note models.Note
	jsonStr := []byte(`{"title":"` + title + `", "content":"` + content + `"}`)
	body, err := sendRequest("PUT", "/api/notes/"+id, jsonStr, token)
	if err != nil {
		return note, err
	}
	// Parse json body
	err = json.Unmarshal(body, &note)
	if err != nil {
		return note, err
	}
	return note, nil
}

// DeleteNote deletes a note
func DeleteNote(id string, token string) (models.Note, error) {
	var note models.Note
	jsonStr := []byte(`{"id":"` + id + `"}`)
	body, err := sendRequest("DELETE", "/api/notes/"+id, jsonStr, token)
	if err != nil {
		return note, err
	}
	// Parse json body
	err = json.Unmarshal(body, &note)
	if err != nil {
		return note, err
	}
	return note, nil
}

func PrintNote(note models.Note) {
	var printableData string
	printableData += "[" + strconv.Itoa(int(note.ID)) + "]" + "\t" + "Account: " + note.Username + "\n"
	printableData += "Title: " + note.Title + "\n"
	printableData += ">>\n" + note.Content + "\n"
	printableData += "Created on: " + note.CreatedAt.String() + "\n"
	printableData += "Updated on: " + note.UpdatedAt.String() + "\n"
	fmt.Println(printableData)
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

// HashAndSalt generates a hashed password
func HashAndSalt(password string) (string, error) {
	// Generate a hashed password with bcypt
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
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
	if u.Username == "" {
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
	if len(s) >= 8 {
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

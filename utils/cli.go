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
	"time"

	"github.com/mrinjamul/gnote/models"
)

// var ApiURL = "https://note.mrinjamul.in"
// var ApiURL = "http://localhost:8080"
var ApiURL = "https://gnote.up.railway.app"

// Response is the response from the API
type Response struct {
	Status  string        `json:"status"`
	Message string        `json:"message"`
	Note    models.Note   `json:"note"`
	Notes   []models.Note `json:"notes"`
	Error   string        `json:"error"`
}

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
	var resp Response
	jsonStr := []byte(`{"title":"` + title + `", "content":"` + content + `"}`)
	body, err := sendRequest("POST", "/api/notes", jsonStr, token)
	if err != nil {
		return models.Note{}, err
	}
	// Parse json body
	err = json.Unmarshal(body, &resp)
	if err != nil {
		return models.Note{}, err
	}
	if resp.Status == "success" || resp.Message == "success" {
		return resp.Note, nil
	}
	return models.Note{}, errors.New(resp.Error)
}

// GetNotes gets all notes
func GetNotes(token string) ([]models.Note, error) {
	var resp Response
	body, err := sendRequest("GET", "/api/notes", nil, token)
	if err != nil {
		return []models.Note{}, err
	}
	// Parse json body
	err = json.Unmarshal(body, &resp)
	if err != nil {
		return []models.Note{}, err
	}
	if resp.Status == "success" || resp.Message == "success" {
		return resp.Notes, nil
	}
	return []models.Note{}, errors.New(resp.Error)
}

// GetNote gets a note
func GetNote(id string, token string) (models.Note, error) {
	var resp Response
	jsonStr := []byte(`{"id":"` + id + `"}`)
	body, err := sendRequest("GET", "/api/notes/"+id, jsonStr, token)
	if err != nil {
		return models.Note{}, err
	}
	// Parse json body
	err = json.Unmarshal(body, &resp)
	if err != nil {
		return models.Note{}, err
	}
	if resp.Status == "success" || resp.Message == "success" {
		return resp.Note, nil
	}
	return models.Note{}, errors.New(resp.Error)
}

// UpdateNote updates a note
func UpdateNote(id, title, content string, token string) (models.Note, error) {
	var resp Response
	jsonStr := []byte(`{"title":"` + title + `", "content":"` + content + `"}`)
	body, err := sendRequest("PUT", "/api/notes/"+id, jsonStr, token)
	if err != nil {
		return models.Note{}, err
	}
	// Parse json body
	err = json.Unmarshal(body, &resp)
	if err != nil {
		return models.Note{}, err
	}
	if resp.Status == "success" || resp.Message == "success" {
		return resp.Note, nil
	}
	return models.Note{}, errors.New(resp.Error)
}

// DeleteNote deletes a note
func DeleteNote(id string, token string) (models.Note, error) {
	var resp Response
	jsonStr := []byte(`{"id":"` + id + `"}`)
	body, err := sendRequest("DELETE", "/api/notes/"+id, jsonStr, token)
	if err != nil {
		return models.Note{}, err
	}
	// Parse json body
	err = json.Unmarshal(body, &resp)
	if err != nil {
		return models.Note{}, err
	}
	if resp.Status == "success" || resp.Message == "success" {
		return resp.Note, nil
	}
	return models.Note{}, errors.New(resp.Error)
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

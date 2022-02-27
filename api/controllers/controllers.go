package controllers

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mrinjamul/gnote/models"
	"github.com/mrinjamul/gnote/repository"
)

type Note interface {
	Create(ctx *gin.Context)
	Read(ctx *gin.Context)
	ReadAll(ctx *gin.Context)
	Update(ctx *gin.Context)
	Delete(ctx *gin.Context)
}

type note struct {
	noteRepo repository.NoteRepo
}

// Create creates a new note
func (n *note) Create(ctx *gin.Context) {
	bytes, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		log.Fatal(err)
	}
	var note models.Note
	err = json.Unmarshal(bytes, &note)
	if err != nil {
		log.Fatal(err)
	}

	err = n.noteRepo.Create(ctx, &note)
	if err != nil {
		ctx.JSON(200, gin.H{
			"message": err,
			"note":    note,
		})
	}
	ctx.JSON(200, gin.H{
		"message": "Success",
		"note":    note,
	})
}

// Read reads a note
func (n *note) Read(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		log.Fatal(err)
	}

	note := models.Note{
		ID: uint64(id),
	}

	note, err = n.noteRepo.Read(ctx, note)
	if err != nil {
		ctx.JSON(
			http.StatusOK,
			gin.H{
				"message": err,
			},
		)
	}
	ctx.JSON(
		http.StatusOK,
		gin.H{
			"message": "success",
			"note":    note,
		},
	)
}

// ReadAll reads all notes
func (n *note) ReadAll(ctx *gin.Context) {
	notes, err := n.noteRepo.ReadAll(ctx)
	if err != nil {
		ctx.JSON(
			http.StatusOK,
			gin.H{
				"message": err,
				"notes":   notes,
			},
		)
	}
	ctx.JSON(
		http.StatusOK,
		gin.H{
			"message": "success",
			"notes":   notes,
		},
	)
}

// Update updates a note
func (n *note) Update(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		log.Fatal(err)
	}

	bytes, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		log.Fatal(err)
	}

	var note, existingNote models.Note
	err = json.Unmarshal(bytes, &note)
	if err != nil {
		log.Fatal(err)
	}

	// get existing note

	if note.Title == "" {
		existingNote.Title = note.Title
	}
	if note.Content == "" {
		existingNote.Content = note.Content
	}
	if note.User == "" {
		existingNote.User = note.User
	}
	// if note.Archived != existingNote.Archived {
	// 	existingNote.Archived = note.Archived
	// }

	note.ID = uint64(id)

	note, err = n.noteRepo.Update(ctx, existingNote)
	if err != nil {
		ctx.JSON(
			http.StatusOK,
			gin.H{
				"message": err,
				"note":    note,
			},
		)
	}
	ctx.JSON(
		http.StatusOK,
		gin.H{
			"message": "success",
			"note":    note,
		},
	)
}

// Delete deletes a note
func (n *note) Delete(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		log.Fatal(err)
	}

	note := models.Note{
		ID: uint64(id),
	}

	err = n.noteRepo.Delete(ctx, &note)
	if err != nil {
		ctx.JSON(
			http.StatusOK,
			gin.H{
				"message": err,
				"note":    note,
			},
		)
	}

	ctx.JSON(
		http.StatusOK,
		gin.H{
			"message": "success",
			"note":    note,
		},
	)
}

// NewNote initializes note
func NewNote(noteRepo repository.NoteRepo) Note {
	return &note{
		noteRepo: noteRepo,
	}
}

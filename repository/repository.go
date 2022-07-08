package repository

import (
	"github.com/gin-gonic/gin"
	"github.com/mrinjamul/gnote/models"
	"gorm.io/gorm"
)

type NoteRepo interface {
	Create(ctx *gin.Context, note *models.Note) error
	Read(ctx *gin.Context, note models.Note) (models.Note, error)
	ReadByUserName(ctx *gin.Context, user string) ([]models.Note, error)
	ReadAll(ctx *gin.Context) ([]models.Note, error)
	Update(ctx *gin.Context, note models.Note) (models.Note, error)
	Delete(ctx *gin.Context, note *models.Note) error
}

type noteRepo struct {
	db gorm.DB
}

// Create creates a new note
func (repo *noteRepo) Create(ctx *gin.Context, note *models.Note) error {
	result := repo.db.Create(note)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// Read reads a note
func (repo *noteRepo) Read(ctx *gin.Context, note models.Note) (models.Note, error) {
	result := repo.db.First(&note)
	if result.Error != nil {
		return note, result.Error
	}
	return note, nil
}

// ReadByUserName reads all notes by user name
func (repo *noteRepo) ReadByUserName(ctx *gin.Context, user string) ([]models.Note, error) {
	var notes []models.Note
	result := repo.db.Find(&notes, "user_name = ?", user)
	if result.Error != nil {
		return notes, result.Error
	}
	return notes, nil
}

// ReadAll reads all notes
func (repo *noteRepo) ReadAll(ctx *gin.Context) ([]models.Note, error) {
	var notes []models.Note
	result := repo.db.Find(&notes)
	if result.Error != nil {
		return notes, result.Error
	}
	return notes, nil
}

// Update updates a note
func (repo *noteRepo) Update(ctx *gin.Context, note models.Note) (models.Note, error) {
	var existingNote models.Note
	result := repo.db.Find(&existingNote, "id = ?", note.ID)
	if result.Error != nil {
		return existingNote, result.Error
	}

	if note.Archived {
		existingNote.Archived = true
	}
	if note.Title != "" {
		existingNote.Title = note.Title
	}
	if note.Content != "" {
		existingNote.Content = note.Content
	}

	result = repo.db.Save(existingNote)
	if result.Error != nil {
		return existingNote, result.Error
	}
	return existingNote, nil
}

// Delete deletes a note
func (repo *noteRepo) Delete(ctx *gin.Context, note *models.Note) error {
	result := repo.db.Delete(note)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// NewNoteRepo initializes the note repository
func NewNoteRepo(db *gorm.DB) NoteRepo {
	return &noteRepo{
		db: *db,
	}
}

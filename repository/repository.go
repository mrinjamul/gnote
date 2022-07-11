package repository

import (
	"github.com/gin-gonic/gin"
	"github.com/mrinjamul/gnote/models"
	"github.com/mrinjamul/gnote/utils"
	"gorm.io/gorm"
)

type NoteRepo interface {
	Create(ctx *gin.Context, note *models.Note) error
	Read(ctx *gin.Context, note *models.Note) error
	ReadByUserName(ctx *gin.Context, user string) ([]models.Note, error)
	ReadAll(ctx *gin.Context) ([]models.Note, error)
	Update(ctx *gin.Context, note models.Note) (models.Note, error)
	Delete(ctx *gin.Context, note *models.Note) error
	DeleteAllByUserName(ctx *gin.Context, username string) error
	VerifyPassword(ctx *gin.Context, username, password string) (bool, error)
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
func (repo *noteRepo) Read(ctx *gin.Context, note *models.Note) error {
	result := repo.db.First(&note)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// ReadByUserName reads all notes by user name
func (repo *noteRepo) ReadByUserName(ctx *gin.Context, user string) ([]models.Note, error) {
	var notes []models.Note
	result := repo.db.Find(&notes, "username = ?", user)
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
	// Save notes
	err := repo.db.Save(&note).Error
	if err != nil {
		return note, err
	}
	return note, nil
}

// Delete deletes a note
func (repo *noteRepo) Delete(ctx *gin.Context, note *models.Note) error {
	result := repo.db.Delete(note)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// 	DeleteAllByUserName deletes all notes by user name
func (repo noteRepo) DeleteAllByUserName(ctx *gin.Context, username string) error {
	var notes []models.Note
	notes, err := repo.ReadByUserName(ctx, username)
	if err != nil {
		return err
	}
	for _, note := range notes {
		err = repo.db.Delete(&note).Error
		if err != nil {
			return err
		}
	}
	return nil
}

// VerifyPassword verifies the password
func (repo *noteRepo) VerifyPassword(ctx *gin.Context, username, password string) (bool, error) {
	var user models.User
	err := repo.db.Find(&user, "username = ?", username).Error
	ok := utils.VerifyHash(password, user.Password)
	if err != nil {
		return false, err
	}
	return ok, nil
}

// NewNoteRepo initializes the note repository
func NewNoteRepo(db *gorm.DB) NoteRepo {
	return &noteRepo{
		db: *db,
	}
}

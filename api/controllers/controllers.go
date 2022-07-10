package controllers

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/mrinjamul/gnote/models"
	"github.com/mrinjamul/gnote/repository"
	"github.com/mrinjamul/gnote/utils"
)

type Note interface {
	Create(ctx *gin.Context)
	Read(ctx *gin.Context)
	ReadAll(ctx *gin.Context)
	Update(ctx *gin.Context)
	Delete(ctx *gin.Context)
	DeleteByUsername(ctx *gin.Context)
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
	// Get cookie "token"
	tokenString, err := ctx.Cookie("token")
	if err != nil {
		tkn, err := utils.ParseToken(ctx.Request.Header.Get("Authorization"))
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": "invalid token",
			})
			ctx.Abort()
			return
		}
		tokenString = tkn
	}

	claims := &models.Claims{}
	token, _ := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtKey), nil
	})
	if !token.Valid {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid token",
		})
		ctx.Abort()
		return
	}

	note.Username = claims.Username

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

	// Get cookie "token"
	tokenString, err := ctx.Cookie("token")
	if err != nil {
		tkn, err := utils.ParseToken(ctx.Request.Header.Get("Authorization"))
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": "invalid token",
			})
			ctx.Abort()
			return
		}
		tokenString = tkn
	}

	claims := &models.Claims{}
	token, _ := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtKey), nil
	})
	if !token.Valid {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid token",
		})
		ctx.Abort()
		return
	}

	note := models.Note{
		ID:       uint64(id),
		Username: claims.Username,
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
	// Get cookie "token"
	tokenString, err := ctx.Cookie("token")
	if err != nil {
		tkn, err := utils.ParseToken(ctx.Request.Header.Get("Authorization"))
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": "invalid token",
			})
			ctx.Abort()
			return
		}
		tokenString = tkn
	}

	claims := &models.Claims{}
	token, _ := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtKey), nil
	})
	if !token.Valid {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid token",
		})
		ctx.Abort()
		return
	}

	notes, err := n.noteRepo.ReadByUserName(ctx, claims.Username)
	if err != nil {
		ctx.JSON(
			http.StatusOK,
			gin.H{
				"message": err,
				"notes":   "",
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

	var note models.Note
	err = json.Unmarshal(bytes, &note)
	if err != nil {
		log.Fatal(err)
	}

	// Get existing note
	var existingNote models.Note
	existingNote = models.Note{
		ID: uint64(id),
	}
	// get existing note
	existingNote, err = n.noteRepo.Read(ctx, existingNote)
	if err != nil {
		log.Fatal(err)
	}

	if note.Title == "" {
		existingNote.Title = note.Title
	}
	if note.Content == "" {
		existingNote.Content = note.Content
	}
	if note.Username == "" {
		existingNote.Username = note.Username
	}
	if note.Archived != existingNote.Archived {
		existingNote.Archived = note.Archived
	}

	// Get cookie "token"
	tokenString, err := ctx.Cookie("token")
	if err != nil {
		tkn, err := utils.ParseToken(ctx.Request.Header.Get("Authorization"))
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": "invalid token",
			})
			ctx.Abort()
			return
		}
		tokenString = tkn
	}

	claims := &models.Claims{}
	token, _ := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtKey), nil
	})
	if !token.Valid {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid token",
		})
		ctx.Abort()
		return
	}

	note.ID = uint64(id)

	// if username is not same as login
	if claims.Username != existingNote.Username {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error":      "you are not owner of this note",
			"pre_notes":  existingNote,
			"claim_user": claims.Username,
			"notes":      note,
		})
		ctx.Abort()
		return
	}

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

	// Get cookie "token"
	tokenString, err := ctx.Cookie("token")
	if err != nil {
		tkn, err := utils.ParseToken(ctx.Request.Header.Get("Authorization"))
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": "invalid token",
			})
			ctx.Abort()
			return
		}
		tokenString = tkn
	}

	claims := &models.Claims{}
	token, _ := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtKey), nil
	})
	if !token.Valid {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid token",
		})
		ctx.Abort()
		return
	}

	note := models.Note{
		ID:       uint64(id),
		Username: claims.Username,
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

// DeleteByUsername deletes all notes by username
func (n *note) DeleteByUsername(ctx *gin.Context) {
	var user map[string]string
	err := ctx.BindJSON(&user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})
		ctx.Abort()
		return
	}
	// Get cookie "token"
	tokenString, err := ctx.Cookie("token")
	if err != nil {
		tkn, err := utils.ParseToken(ctx.Request.Header.Get("Authorization"))
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": "invalid token",
			})
			ctx.Abort()
			return
		}
		tokenString = tkn
	}

	claims := &models.Claims{}
	token, _ := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtKey), nil
	})
	if !token.Valid {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid token",
		})
		ctx.Abort()
		return
	}

	valid, err := n.noteRepo.VerifyPassword(ctx, claims.Username, user["password"])
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})
		ctx.Abort()
		return
	}

	if !valid {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid password",
		})
		ctx.Abort()
		return
	}

	err = n.noteRepo.DeleteAllByUserName(ctx, claims.Username)
	if err != nil {
		ctx.JSON(
			http.StatusOK,
			gin.H{
				"status":  "error",
				"message": err,
			})
	}
	ctx.JSON(
		http.StatusOK,
		gin.H{
			"status":  "success",
			"message": "all notes deleted",
		})
}

// NewNote initializes note
func NewNote(noteRepo repository.NoteRepo) Note {
	return &note{
		noteRepo: noteRepo,
	}
}

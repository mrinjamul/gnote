package handlers

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mrinjamul/gnote/database"
	"github.com/mrinjamul/gnote/models"
	"gorm.io/gorm"
)

var db *gorm.DB

func init() {
	db = database.GetDB()
}

func GetNotes(ctx *gin.Context) {
	var notes []models.Note
	result := db.Find(&notes)
	if result.Error != nil {
		ctx.JSON(
			http.StatusOK,
			gin.H{
				"message": result.Error,
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

func GetNote(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		log.Fatal(err)
	}

	note := models.Note{
		ID: uint64(id),
	}

	result := db.First(&note)
	if result.Error != nil {
		ctx.JSON(
			http.StatusOK,
			gin.H{
				"message": result.Error,
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

func CreateNote(ctx *gin.Context) {
	bytes, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		log.Fatal(err)
	}
	var note models.Note
	err = json.Unmarshal(bytes, &note)
	if err != nil {
		log.Fatal(err)
	}
	result := db.Omit("id").Create(&note)
	if result.Error != nil {
		ctx.JSON(
			http.StatusExpectationFailed,
			note,
		)
	}

	// Create a new note
	ctx.JSON(
		http.StatusOK,
		note,
	)
}

func UpdateNote(ctx *gin.Context) {
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

	note.ID = uint64(id)

	result := db.Save(&note)
	if result.Error != nil {
		ctx.JSON(
			http.StatusOK,
			gin.H{
				"message": result.Error,
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

func DeleteNote(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		log.Fatal(err)
	}

	note := models.Note{
		ID: uint64(id),
	}

	result := db.Delete(&note)
	if result.Error != nil {
		ctx.JSON(
			http.StatusOK,
			gin.H{
				"message": result.Error,
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

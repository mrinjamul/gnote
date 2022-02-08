package database

import (
	"github.com/mrinjamul/gnote/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func GetDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("sqlite.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	// Migrate the schema
	db.AutoMigrate(&models.Note{})
	return db
}

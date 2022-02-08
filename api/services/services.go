package services

import (
	"github.com/mrinjamul/gnote/api/controllers"
	"github.com/mrinjamul/gnote/database"
	"github.com/mrinjamul/gnote/repository"
)

type Services interface {
	NoteService() controllers.Note
}

type services struct {
	note controllers.Note
}

func (svc *services) NoteService() controllers.Note {
	return svc.note
}

// NewServices initializes services
func NewServices() Services {
	db := database.GetDB()
	return &services{
		note: controllers.NewNote(
			repository.NewNoteRepo(db),
		),
	}
}

package services

import (
	"github.com/mrinjamul/gnote/api/controllers"
	"github.com/mrinjamul/gnote/database"
	"github.com/mrinjamul/gnote/repository"
)

type Services interface {
	HealthCheckService() controllers.HealthCheck
	NoteService() controllers.Note
}

type services struct {
	healthCheck controllers.HealthCheck
	note        controllers.Note
}

func (svc *services) HealthCheckService() controllers.HealthCheck {
	return svc.healthCheck
}

func (svc *services) NoteService() controllers.Note {
	return svc.note
}

// NewServices initializes services
func NewServices() Services {
	db := database.GetDB()
	return &services{
		healthCheck: controllers.NewHealthCheck(),
		note: controllers.NewNote(
			repository.NewNoteRepo(db),
		),
	}
}

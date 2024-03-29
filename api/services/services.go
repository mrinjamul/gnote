package services

import (
	"github.com/mrinjamul/gnote/api/controllers"
	"github.com/mrinjamul/gnote/database"
	"github.com/mrinjamul/gnote/repository"
)

type Services interface {
	HealthCheckService() controllers.HealthCheck
	NoteService() controllers.Note
	UserService() controllers.User
	ViewService() controllers.Views
}

type services struct {
	healthCheck controllers.HealthCheck
	note        controllers.Note
	user        controllers.User
	views       controllers.Views
}

func (svc *services) HealthCheckService() controllers.HealthCheck {
	return svc.healthCheck
}

func (svc *services) NoteService() controllers.Note {
	return svc.note
}

func (svc *services) UserService() controllers.User {
	return svc.user
}

func (svc *services) ViewService() controllers.Views {
	return svc.views
}

// NewServices initializes services
func NewServices() Services {
	db := database.GetDB()
	return &services{
		healthCheck: controllers.NewHealthCheck(),
		note: controllers.NewNote(
			repository.NewNoteRepo(db),
		),
		user: controllers.NewUser(
			repository.NewUserRepo(db),
		),
		views: controllers.NewViews(),
	}
}

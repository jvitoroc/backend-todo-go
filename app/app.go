package app

import (
	"github.com/jvitoroc/todo-go/auth"
	"github.com/jvitoroc/todo-go/config"
	"github.com/jvitoroc/todo-go/email"
	"github.com/jvitoroc/todo-go/repository"
)

type App struct {
	Repository   *repository.Repository
	EmailService *email.EmailService
	AuthService  *auth.AuthService
	Config       *config.Config
}

func NewApp(repo *repository.Repository, email *email.EmailService, auth *auth.AuthService, config *config.Config) *App {
	return &App{
		Repository:   repo,
		EmailService: email,
		AuthService:  auth,
		Config:       config,
	}
}

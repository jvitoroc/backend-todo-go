package api

import (
	"github.com/gorilla/mux"
	"github.com/jvitoroc/todo-go/app"
)

type API struct {
	App    *app.App
	Router *Router

	MainRouter *mux.Router
}

func NewAPI(app *app.App, router *mux.Router) *API {
	api := &API{
		App:        app,
		MainRouter: router,
		Router:     &Router{},
	}

	api.setupRoutes()

	return api
}

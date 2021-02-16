package api

import (
	"github.com/gorilla/mux"
)

type Router struct {
	Root                *mux.Router
	User                *mux.Router
	Session             *mux.Router
	VerificationRequest *mux.Router
	Todo                *mux.Router
}

func (api *API) setupRoutes() {
	api.Router.Root = api.MainRouter
	api.Router.User = api.MainRouter.PathPrefix("/user").Subrouter()
	api.Router.Session = api.Router.User.PathPrefix("/session").Subrouter()
	api.Router.VerificationRequest = api.Router.User.PathPrefix("/verification-request").Subrouter()
	api.Router.Todo = api.MainRouter.PathPrefix("/todo").Subrouter()

	api.InitUser()
	api.InitSession()
	api.InitVerificationRequest()
	api.InitTodo()
}

package api

import (
	"net/http"

	hn "github.com/jvitoroc/todo-go/api/handler"
	"github.com/jvitoroc/todo-go/model"
)

func (api *API) InitSession() {
	api.Router.Session.Handle("", api.createHandler(api.CreateSession)).Methods("POST")
	api.Router.Session.Handle("/google", api.createHandler(api.CreateGoogleSession)).Methods("POST")
}

func (api *API) CreateSession(ctx *hn.RequestContext, w http.ResponseWriter, r *http.Request) hn.Response {
	bs, err := model.BasicSessionFromSession(r.Body)
	if err != nil {
		return err
	}

	if err := bs.Validate(); err != nil {
		return err
	}

	token, err := api.App.CreateSession(bs)
	if err != nil {
		return err
	}

	return model.NewCreatedResponse(model.MSG_SESSION_CREATED).AddObject("token", token)
}

func (api *API) CreateGoogleSession(ctx *hn.RequestContext, w http.ResponseWriter, r *http.Request) hn.Response {
	gs, err := model.GoogleSessionFromSession(r.Body)
	if err != nil {
		return err
	}

	if err := gs.Validate(); err != nil {
		return err
	}

	token, err := api.App.CreateGoogleSession(gs)
	if err != nil {
		return err
	}

	return model.NewCreatedResponse(model.MSG_SESSION_CREATED).AddObject("token", token)
}

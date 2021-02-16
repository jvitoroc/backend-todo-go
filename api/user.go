package api

import (
	"net/http"

	hn "github.com/jvitoroc/todo-go/api/handler"
	"github.com/jvitoroc/todo-go/model"
)

func (api *API) InitUser() {
	api.Router.User.Handle("", api.createHandler(api.CreateUser)).Methods("POST")
	api.Router.User.Handle("", api.createProtectedHandler(api.GetUser, true)).Methods("GET")
}

func (api *API) CreateUser(ctx *hn.RequestContext, w http.ResponseWriter, r *http.Request) hn.Response {
	user, err := model.UserFromJson(r.Body)
	if err != nil {
		return err
	}

	if err := user.Validate(); err != nil {
		return err
	}

	user.Verified = false
	if err := api.App.CreateUser(user); err != nil {
		return err
	}

	user.OmitSecretFields()

	return model.NewCreatedResponse(model.MSG_USER_CREATED).AddObject("user", user)
}

func (api *API) GetUser(ctx *hn.RequestContext, w http.ResponseWriter, r *http.Request) hn.Response {
	return model.NewOKResponse(model.MSG_USER_RETRIEVED).AddObject("user", ctx.CurrentUser)
}

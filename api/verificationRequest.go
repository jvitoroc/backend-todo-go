package api

import (
	"net/http"

	hn "github.com/jvitoroc/todo-go/api/handler"
	"github.com/jvitoroc/todo-go/model"
)

func (api *API) InitVerificationRequest() {
	api.Router.VerificationRequest.Handle("", api.createProtectedHandler(api.CheckVerificationRequest, false)).Methods("POST")
	api.Router.VerificationRequest.Handle("", api.createProtectedHandler(api.UpsertVerificationRequest, false)).Methods("PUT")
}

func (api *API) CheckVerificationRequest(ctx *hn.RequestContext, w http.ResponseWriter, r *http.Request) hn.Response {
	if ctx.CurrentUser.Verified {
		return model.NewBadRequestError(model.MSG_USER_ALREADY_VERIFIED).AddObject("user", ctx.CurrentUser)
	}

	verify, err := model.VerifyUserAccountFromJson(r.Body)
	if err != nil {
		return err
	}

	if err := verify.Validate(); err != nil {
		return err
	}

	if err := api.App.CheckVerificationRequest(ctx.CurrentUser, verify.VerificationCode); err != nil {
		return err
	}

	return model.NewOKResponse(model.MSG_USER_VERIFIED)
}

func (api *API) UpsertVerificationRequest(ctx *hn.RequestContext, w http.ResponseWriter, r *http.Request) hn.Response {
	if ctx.CurrentUser.Verified {
		return model.NewBadRequestError(model.MSG_USER_ALREADY_VERIFIED).AddObject("user", ctx.CurrentUser)
	}

	if err := api.App.UpsertVerificationRequest(ctx.CurrentUser); err != nil {
		return err
	}

	return model.NewOKResponse(model.MSG_VERIFICATION_SENT)
}

package api

import hn "github.com/jvitoroc/todo-go/api/handler"

func (api *API) createHandler(handler hn.HandlerFunc) *hn.Handler {
	return &hn.Handler{
		App:               api.App,
		Authenticate:      false,
		CheckVerification: false,
		Handler:           handler,
	}
}

func (api *API) createProtectedHandler(handler hn.HandlerFunc, checkVerification bool) *hn.Handler {
	return &hn.Handler{
		App:               api.App,
		Authenticate:      true,
		CheckVerification: checkVerification,
		Handler:           handler,
	}
}

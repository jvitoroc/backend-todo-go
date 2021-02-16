package handler

import (
	"net/http"

	"github.com/jvitoroc/todo-go/app"
	"github.com/jvitoroc/todo-go/model"
)

type Handler struct {
	Authenticate      bool
	CheckVerification bool

	App *app.App

	Handler HandlerFunc
}

type RequestContext struct {
	CurrentUser *model.User

	Request *http.Request
}

type HandlerFunc func(*RequestContext, http.ResponseWriter, *http.Request) Response

func (hn *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := &RequestContext{
		Request: r,
	}

	if err := hn.RunAllMiddlewares(ctx); err != nil {
		writeResponse(w, err)
		return
	}

	res := hn.Handler(ctx, w, r)
	writeResponse(w, res)
}

func writeResponse(w http.ResponseWriter, res Response) {
	w.WriteHeader(res.GetCode())
	w.Write([]byte(res.ToJson()))
}

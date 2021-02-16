package handler

import (
	"strings"

	"github.com/jvitoroc/todo-go/model"
)

func (h *Handler) RunAllMiddlewares(ctx *RequestContext) *model.AppError {
	if h.Authenticate {
		if err := h.CheckSessionMiddleware(ctx); err != nil {
			return err
		}
	}

	if h.CheckVerification {
		if err := h.CheckVerificationMiddleware(ctx); err != nil {
			return err
		}
	}

	return nil
}

func (h *Handler) CheckSessionMiddleware(ctx *RequestContext) *model.AppError {
	token, err := h.getToken(ctx)
	if err != nil {
		return err
	}

	userId, err := h.App.VerifyToken(token)
	if err != nil {
		return err
	}

	user, err := h.App.GetUser(userId)
	if err != nil {
		return err
	}

	user.OmitSecretFields()
	ctx.CurrentUser = user
	return nil
}

func (h *Handler) CheckVerificationMiddleware(ctx *RequestContext) *model.AppError {
	if !ctx.CurrentUser.Verified {
		return model.NewForbiddenError("User account not yet verified.").AddObject("user", ctx.CurrentUser)
	}

	return nil
}

func (h *Handler) getToken(ctx *RequestContext) (string, *model.AppError) {
	token := ctx.Request.Header.Get("Authorization")
	if strings.HasPrefix(token, "Bearer ") == false {
		return "", model.NewBadRequestError("Authorization token not provided.")
	}

	return token[7:], nil
}

package app

import (
	"net/http"

	"github.com/jvitoroc/todo-go/model"
	"github.com/jvitoroc/todo-go/util"
)

func (app *App) CreateSession(bs *model.BasicSession) (string, *model.AppError) {
	user, err := app.Repository.GetUserByUsername(bs.Username)
	if err != nil {
		if err.Code == http.StatusNotFound {
			err = model.NewBadRequestError(model.MSG_INVALID_CREDENTIAL)
		}
		return "", err
	}

	if !util.VerifyPasswordHash(user.Password, bs.Password) {
		return "", model.NewBadRequestError(model.MSG_INVALID_CREDENTIAL)
	}

	token, err := app.SignToken(user.ID)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (app *App) CreateGoogleSession(bs *model.GoogleSession) (string, *model.AppError) {
	if app.AuthService.Google == nil {
		return "", model.NewInternalError(model.MSG_GOOGLE_UNAVAILABLE)
	}

	googleRes, err := app.AuthService.Google.Validate(*bs.IdToken)
	if err != nil {
		return "", model.NewInternalError(model.MSG_GOOGLE_TOKEN_ERROR).SetDetail(err.Error())
	}

	var user *model.User
	user, appErr := app.Repository.GetUserByGoogleSub(googleRes.Subject)
	if appErr != nil && appErr.Code == http.StatusNotFound {
		user, appErr = app.CreateUserWithGoogleClaims(googleRes)
		if appErr != nil {
			return "", appErr
		}
	} else if appErr != nil {
		return "", appErr
	}

	token, signErr := app.SignToken(user.ID)
	if signErr != nil {
		return "", signErr
	}

	return token, nil
}

package app

import (
	"github.com/jvitoroc/todo-go/model"
	"github.com/jvitoroc/todo-go/repository"
	"github.com/jvitoroc/todo-go/util"
	"google.golang.org/api/idtoken"
)

func (app *App) CreateUser(user *model.User) *model.AppError {
	var passwordHash string
	if err := util.GeneratePasswordHash(&passwordHash, user.Password); err != nil {
		return model.NewGenericInternalError(err)
	}

	user.Password = passwordHash

	return app.Repository.BeginTran(func(tran *repository.Repository) *model.AppError {
		usernameExists, err := tran.CheckIfUsernameExists(user.Username)
		if err != nil {
			return err
		}

		emailExists, err := tran.CheckIfEmailExists(user.Email)
		if err != nil {
			return err
		}

		if usernameExists || emailExists {
			err := model.NewFormError(nil)
			if usernameExists {
				err.AddError("username", model.MSG_USER_USERNAME_ALREADY_EXISTS)
			}
			if emailExists {
				err.AddError("email", model.MSG_USER_EMAIL_ALREADY_EXISTS)
			}
			return err
		}

		user, err = tran.CreateUser(user)
		if err != nil {
			return err
		}

		return nil
	})
}

func (app *App) CreateUserWithGoogleClaims(payload *idtoken.Payload) (*model.User, *model.AppError) {
	email := payload.Claims["email"].(string)

	user := model.User{
		Username:  email,
		Email:     email,
		GoogleSub: &payload.Subject,
		Verified:  true,
	}

	err := app.CreateUser(&user)
	if err != nil {
		return nil, err
	}

	return &user, err
}

func (app *App) GetUser(userId int) (*model.User, *model.AppError) {
	return app.Repository.GetUser(userId)
}

func (app *App) VerifyUser(userId int) *model.AppError {
	user, err := app.Repository.GetUser(userId)
	if err != nil {
		return err
	}

	user.Verified = true

	return app.Repository.UpdateUser(user)
}

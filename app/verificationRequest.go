package app

import (
	"math/rand"
	"time"

	"github.com/jvitoroc/todo-go/model"
)

func (app *App) CheckVerificationRequest(user *model.User, verificationCode string) *model.AppError {
	vr, err := app.Repository.GetVerificationRequest(user.ID)
	if err != nil {
		return err
	}

	if verificationCode != vr.Code || time.Now().After(vr.ExpiresAt) {
		return model.NewBadRequestError(model.MSG_INVALID_CODE)
	}

	if err != app.VerifyUser(user.ID) {
		return err
	}

	return nil
}

func (app *App) UpsertVerificationRequest(user *model.User) *model.AppError {
	vr := &model.VerificationRequest{
		UserID:    user.ID,
		Code:      app.GenerateVerificationCode(),
		ExpiresAt: time.Now().Local().Add(time.Minute * time.Duration(model.VERIFICATION_CODE_EXPIRATION)),
	}

	dbVr, err := app.Repository.UpsertVerificationRequest(vr)
	if err != nil {
		return err
	}

	return app.SendVerificationEmail(user, dbVr)
}

func (app *App) SendVerificationEmail(user *model.User, vr *model.VerificationRequest) *model.AppError {
	body :=
		"Hey " + user.Username + ", here's the code needed to activate your account: " + vr.Code + "\r\n" +
			"It will expire on " + vr.ExpiresAt.Format(time.RFC1123)

	if err := app.EmailService.SendEmail(user.Email, "Todo App: here is your verification code", body); err != nil {
		return model.NewGenericInternalError(err)
	}

	return nil
}

func (app *App) GenerateVerificationCode() string {
	code := ""
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	chars := model.VERIFICATION_CODE_CHARS
	lastIndex := len(chars) - 1
	for i := 0; i < model.VERIFICATION_CODE_LENGTH; i++ {
		code = code + string(chars[r.Intn(lastIndex)])
	}

	return code
}

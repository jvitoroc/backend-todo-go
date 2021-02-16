package model

import (
	"io"

	"github.com/jvitoroc/todo-go/util"
)

type BasicSession struct {
	Username string `gorm:"uniqueIndex" json:"username"`
	Password string `json:"password,omitempty"`
}

type GoogleSession struct {
	IdToken *string `json:"idToken"`
}

func BasicSessionFromSession(data io.Reader) (*BasicSession, *AppError) {
	bs := &BasicSession{}
	if err := util.FromJson(data, bs); err != nil {
		return nil, NewGenericBadRequestError(err)
	}

	return bs, nil
}

func GoogleSessionFromSession(data io.Reader) (*GoogleSession, *AppError) {
	gs := &GoogleSession{}
	if err := util.FromJson(data, gs); err != nil {
		return nil, NewGenericBadRequestError(err)
	}

	return gs, nil
}

func (bs *BasicSession) Validate() *AppError {
	if len(bs.Username) < USERNAME_MINIMUM_LENGTH || len(bs.Password) < PASSWORD_MINIMUM_LENGTH {
		return NewBadRequestError(MSG_INVALID_CREDENTIAL)
	}

	return nil
}

func (gs *GoogleSession) Validate() *AppError {
	if gs.IdToken == nil {
		return NewBadRequestError(MSG_TOKEN_NOT_PROVIDED)
	}

	return nil
}

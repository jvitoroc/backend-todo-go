package model

import (
	"fmt"
	"io"
	"time"

	"github.com/jvitoroc/todo-go/util"
)

type User struct {
	ID       int    `gorm:"primaryKey;autoIncrement" json:"userId"`
	Username string `gorm:"uniqueIndex" json:"username"`
	Email    string `gorm:"unique" json:"email"`
	Password string `json:"password,omitempty"`
	Verified bool   `json:"verified"`

	GoogleSub *string `gorm:"uniqueIndex" json:"-"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func UserFromJson(data io.Reader) (*User, *AppError) {
	user := &User{}
	if err := util.FromJson(data, user); err != nil {
		return nil, NewGenericBadRequestError(err)
	}

	return user, nil
}

func (user *User) OmitSecretFields() {
	user.Password = ""
}

func (user *User) Validate() *AppError {
	errors := map[string]string{}

	if user.Username == "" {
		errors["username"] = MSG_USER_USERNAME_MISSING
	} else if len(user.Username) < USERNAME_MINIMUM_LENGTH {
		errors["username"] = fmt.Sprintf(MSG_USER_USERNAME_LENGTH, USERNAME_MINIMUM_LENGTH)
	}

	if user.Email == "" {
		errors["email"] = MSG_USER_EMAIL_MISSING
	} else if !util.IsEmailValid(user.Email) {
		errors["email"] = MSG_USER_EMAIL_INVALID
	}

	if user.Password == "" {
		errors["password"] = MSG_USER_PASSWORD_MISSING
	} else if len(user.Password) < PASSWORD_MINIMUM_LENGTH {
		errors["password"] = fmt.Sprintf(MSG_USER_PASSWORD_LENGTH, PASSWORD_MINIMUM_LENGTH)
	}

	if len(errors) == 0 {
		return nil
	} else {
		return NewFormError(errors)
	}
}

func (user *User) ValidateForLogin() *AppError {
	if len(user.Username) < USERNAME_MINIMUM_LENGTH || len(user.Password) < PASSWORD_MINIMUM_LENGTH {
		return NewBadRequestError(MSG_INVALID_CREDENTIAL)
	}

	return nil
}

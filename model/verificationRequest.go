package model

import (
	"io"
	"time"

	"github.com/jvitoroc/todo-go/util"
)

type VerificationRequest struct {
	UserID    int  `gorm:"primaryKey"`
	User      User `gorm:"constraint:OnDelete:CASCADE;foreignkey:UserID;references:ID"`
	Code      string
	ExpiresAt time.Time
}

type VerifyUserAccount struct {
	VerificationCode string `json:"verificationCode"`
}

func VerifyUserAccountFromJson(data io.Reader) (*VerifyUserAccount, *AppError) {
	verify := &VerifyUserAccount{}
	if err := util.FromJson(data, verify); err != nil {
		return nil, NewGenericBadRequestError(err)
	}

	return verify, nil
}

func (verify *VerifyUserAccount) Validate() *AppError {
	errors := map[string]string{}

	if len(verify.VerificationCode) < VERIFICATION_CODE_LENGTH {
		errors["verificationCode"] = MSG_INVALID_CODE
	}

	if len(errors) == 0 {
		return nil
	} else {
		return NewFormError(errors)
	}
}

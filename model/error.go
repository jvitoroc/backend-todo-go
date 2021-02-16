package model

import (
	"encoding/json"
	"net/http"
)

type AppError struct {
	Code int `json:"-"`

	Message string `json:"message"`
	Detail  string `json:"detail,omitempty"`

	Errors map[string]string `json:"errors,omitempty"`

	Data map[string]interface{} `json:"data,omitempty"`
}

func (err *AppError) GetCode() int {
	return err.Code
}

func (err *AppError) ToJson() string {
	json, _ := json.Marshal(err)
	return string(json)
}

func (err *AppError) SetDetail(detail string) *AppError {
	err.Detail = detail
	return err
}

func (err *AppError) AddObject(key string, value interface{}) *AppError {
	if err.Data == nil {
		err.Data = make(map[string]interface{})
	}

	err.Data[key] = value

	return err
}

func (err *AppError) AddError(key, value string) *AppError {
	if err.Errors == nil {
		err.Errors = make(map[string]string)
	}

	err.Errors[key] = value

	return err
}

func NewGenericInternalError(err error) *AppError {
	return &AppError{Code: http.StatusInternalServerError, Message: MSG_ERR_INTERNAL, Detail: err.Error()}
}

func NewGenericBadRequestError(err error) *AppError {
	return &AppError{Code: http.StatusBadRequest, Message: MSG_ERR_INVALID, Detail: err.Error()}
}

func NewInternalError(message string) *AppError {
	return &AppError{Code: http.StatusInternalServerError, Message: message}
}

func NewBadRequestError(message string) *AppError {
	return &AppError{Code: http.StatusBadRequest, Message: message}
}

func NewUnauthorizedError(message string) *AppError {
	return &AppError{Code: http.StatusUnauthorized, Message: message}
}

func NewForbiddenError(message string) *AppError {
	return &AppError{Code: http.StatusUnauthorized, Message: message}
}

func NewConflictError(message string) *AppError {
	return &AppError{Code: http.StatusConflict, Message: message}
}

func NewNotFoundError(message string) *AppError {
	return &AppError{Code: http.StatusNotFound, Message: message}
}

func NewFormError(errors map[string]string) *AppError {
	return &AppError{Code: http.StatusBadRequest, Message: MSG_ERR_SEVERAL, Errors: errors}
}

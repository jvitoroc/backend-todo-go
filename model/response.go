package model

import (
	"encoding/json"
	"net/http"
)

type AppResponse struct {
	Code int `json:"-"`

	Message string `json:"message"`
	Detail  string `json:"detail,omitempty"`

	Data map[string]interface{} `json:"data,omitempty"`
}

func (res *AppResponse) GetCode() int {
	return res.Code
}

func (res *AppResponse) ToJson() string {
	json, _ := json.Marshal(res)
	return string(json)
}

func (res *AppResponse) AddObject(key string, value interface{}) *AppResponse {
	if res.Data == nil {
		res.Data = make(map[string]interface{})
	}

	res.Data[key] = value

	return res
}

func NewCreatedResponse(message string) *AppResponse {
	return &AppResponse{Code: http.StatusCreated, Message: message}
}

func NewOKResponse(message string) *AppResponse {
	return &AppResponse{Code: http.StatusOK, Message: message}
}

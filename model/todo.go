package model

import (
	"io"
	"time"

	"github.com/jvitoroc/todo-go/util"
)

type Todo struct {
	ID           int       `gorm:"primaryKey;autoIncrement" json:"todoId"`
	ParentTodoID *int      `json:"parentTodoId"`
	ParentTodo   *Todo     `gorm:"constraint:OnDelete:CASCADE;foreignkey:ParentTodoID;references:ID" json:"-"`
	UserID       int       `json:"userId"`
	User         User      `gorm:"constraint:OnDelete:CASCADE;foreignkey:UserID;references:ID" json:"-"`
	Description  string    `json:"description"`
	Completed    bool      `json:"completed"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

type UpdateTodo struct {
	ID          int     `json:"todoId"`
	Description *string `json:"description"`
	Completed   *bool   `json:"completed"`
}

type DeleteManyTodos struct {
	IDs []int `json:"ids"`
}

func TodoFromJson(data io.Reader) (*Todo, *AppError) {
	todo := &Todo{}
	if err := util.FromJson(data, todo); err != nil {
		return nil, NewGenericBadRequestError(err)
	}

	return todo, nil
}

func UpdateTodoFromJson(data io.Reader) (*UpdateTodo, *AppError) {
	update := &UpdateTodo{}
	if err := util.FromJson(data, update); err != nil {
		return nil, NewGenericBadRequestError(err)
	}

	return update, nil
}

func DeleteManyTodosFromJson(data io.Reader) (*DeleteManyTodos, *AppError) {
	delete := &DeleteManyTodos{}
	if err := util.FromJson(data, delete); err != nil {
		return nil, NewGenericBadRequestError(err)
	}

	return delete, nil
}

func (todo *Todo) Validate() *AppError {
	errors := map[string]string{}

	if todo.Description == "" {
		errors["description"] = MSG_TODO_DESCRIPTION_MISSING
	}

	if len(errors) == 0 {
		return nil
	} else {
		return NewFormError(errors)
	}
}

func (todo *UpdateTodo) Validate() *AppError {
	errors := map[string]string{}

	if todo.Description != nil && *todo.Description == "" {
		errors["description"] = MSG_TODO_DESCRIPTION_MISSING
	}

	if len(errors) == 0 {
		return nil
	} else {
		return NewFormError(errors)
	}
}

func (todo *DeleteManyTodos) Validate() *AppError {
	errors := map[string]string{}

	if todo.IDs != nil {
		errors["ids"] = MSG_TODO_IDS_NOT_PROVIDED
	}

	if len(errors) == 0 {
		return nil
	} else {
		return NewFormError(errors)
	}
}

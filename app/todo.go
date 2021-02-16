package app

import (
	"github.com/jvitoroc/todo-go/model"
)

func (app *App) CreateTodo(todo *model.Todo) *model.AppError {
	todo, err := app.Repository.CreateTodo(todo)
	if err != nil {
		return err
	}

	return nil
}

func (app *App) GetTodo(todoId, userId int) (*model.Todo, *model.AppError) {
	return app.Repository.GetTodo(todoId, userId)
}

func (app *App) GetTodosBottomToTop(todoId, userId, depth int) ([]model.Todo, *model.AppError) {
	var i int = 0
	list := make([]model.Todo, 0)
	id := todoId

	for i < depth {
		head, err := app.Repository.GetTodo(id, userId)
		if err != nil {
			return nil, err
		}

		list = append(list, *head)

		if head.ParentTodoID == nil {
			break
		}

		id = *head.ParentTodoID
		i += 1
	}

	return list, nil
}

func (app *App) GetTodoChildren(todoId, userId int) ([]model.Todo, *model.AppError) {
	return app.Repository.GetTodoChildren(todoId, userId)
}

func (app *App) GetRootTodoChildren(userId int) ([]model.Todo, *model.AppError) {
	return app.Repository.GetRootTodoChildren(userId)
}

func (app *App) UpdateTodo(todo *model.UpdateTodo, userId int) (*model.Todo, *model.AppError) {
	dbTodo, err := app.GetTodo(todo.ID, userId)
	if err != nil {
		return nil, err
	}

	if todo.Description != nil {
		dbTodo.Description = *todo.Description
	}

	if todo.Completed != nil {
		dbTodo.Completed = *todo.Completed
	}

	if err := app.Repository.UpdateTodo(dbTodo); err != nil {
		return nil, err
	}

	return dbTodo, nil
}

func (app *App) DeleteTodo(todoId, userId int) *model.AppError {
	return app.Repository.DeleteTodo(todoId, userId)
}

func (app *App) DeleteManyTodos(todoIds []int, userId int) *model.AppError {
	return app.Repository.DeleteManyTodos(todoIds, userId)
}

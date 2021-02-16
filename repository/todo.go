package repository

import (
	"errors"
	"fmt"

	"github.com/jvitoroc/todo-go/model"
	"gorm.io/gorm"
)

func (t *Repository) CreateTodo(todo *model.Todo) (*model.Todo, *model.AppError) {
	if err := t.DB.Create(todo).Error; err != nil {
		return nil, model.NewGenericInternalError(err)
	}

	return todo, nil
}

func (t *Repository) GetTodo(todoId int, userId int) (*model.Todo, *model.AppError) {
	todo := model.Todo{}
	if err := t.DB.Where("user_id = ?", userId).First(&todo, todoId).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, model.NewNotFoundError(fmt.Sprintf(model.MSG_TODO_NOT_FOUND, todoId))
		} else {
			return nil, model.NewGenericInternalError(err)
		}
	}

	return &todo, nil
}

func (t *Repository) GetTodoChildren(todoId int, userId int) ([]model.Todo, *model.AppError) {
	todos := []model.Todo{}
	if err := t.DB.Where("parent_todo_id = ? and user_id = ?", todoId, userId).Order("created_at DESC").Find(&todos).Error; err != nil {
		return nil, model.NewGenericInternalError(err)
	}

	return todos, nil
}

func (t *Repository) GetRootTodoChildren(userId int) ([]model.Todo, *model.AppError) {
	todos := []model.Todo{}
	if err := t.DB.Where("parent_todo_id is null and user_id = ?", userId).Order("created_at DESC").Find(&todos).Error; err != nil {
		return nil, model.NewGenericInternalError(err)
	}

	return todos, nil
}

func (t *Repository) UpdateTodo(todo *model.Todo) *model.AppError {
	if err := t.DB.Save(todo).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.NewBadRequestError(fmt.Sprintf(model.MSG_TODO_NOT_FOUND, todo.ID))
		} else {
			return model.NewGenericInternalError(err)
		}
	}

	return nil
}

func (t *Repository) DeleteTodo(todoId, userId int) *model.AppError {
	var result *gorm.DB
	if result = t.DB.Where("id = ? and user_id = ?", todoId, userId).Delete(&model.Todo{}); result.Error != nil {
		return model.NewGenericInternalError(result.Error)
	}

	if result.RowsAffected == 0 {
		return model.NewNotFoundError(fmt.Sprintf(model.MSG_TODO_NOT_FOUND, todoId))
	}

	return nil
}

func (t *Repository) DeleteManyTodos(todoIds []int, userId int) *model.AppError {
	var result *gorm.DB
	if result = t.DB.Where("id in ? and user_id = ?", todoIds, userId).Delete(&model.Todo{}); result.Error != nil {
		return model.NewGenericInternalError(result.Error)
	}

	return nil
}

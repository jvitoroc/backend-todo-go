package repo

import (
	"errors"
	"fmt"
	"time"

	"github.com/jvitoroc/todo-api/resources/common"
	"gorm.io/gorm"
)

type Todo struct {
	ID           int       `json:"todoId"`
	ParentTodoID *int      `json:"parentTodoId"`
	ParentTodo   *Todo     `json:"-" gorm:"constraint:OnDelete:CASCADE;foreignkey:ID;references:ParentTodoID"`
	UserID       int       `json:"userId"`
	User         User      `json:"-" gorm:"constraint:OnDelete:CASCADE;foreignkey:ID;references:UserID"`
	Description  string    `json:"description"`
	Completed    bool      `json:"completed"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

func InsertTodo(db *gorm.DB, userId int, description string) (*Todo, *common.Error) {
	now := time.Now()

	todo := Todo{UserID: userId, Description: description, CreatedAt: now, UpdatedAt: now}
	if result := db.Create(&todo); result.Error != nil {
		return nil, common.CreateGenericInternalError(result.Error)
	}

	return &todo, nil
}

func InsertTodoChild(db *gorm.DB, userId int, parentTodoId int, description string) (*Todo, *common.Error) {
	now := time.Now()

	todo := Todo{UserID: userId, ParentTodoID: &parentTodoId, Description: description, CreatedAt: now, UpdatedAt: now}
	if result := db.Create(&todo); result.Error != nil {
		if result.Error.Error() == "FOREIGN KEY constraint failed" {
			return nil, common.CreateNotFoundError(fmt.Sprintf("Todo not found under given id (%d).", parentTodoId))
		} else {
			return nil, common.CreateGenericInternalError(result.Error)
		}
	}

	return &todo, nil
}

func UpdateTodo(db *gorm.DB, userId int, todoId int, columns map[string]interface{}) *common.Error {
	var result *gorm.DB

	if result = db.Model(&Todo{}).Where("id = ?", todoId).Updates(columns); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return common.CreateNotFoundError(fmt.Sprintf("Todo not found under given id (%d).", todoId))
		} else {
			return common.CreateGenericInternalError(result.Error)
		}
	}

	return nil
}

func DeleteTodo(db *gorm.DB, userId int, todoId int) *common.Error {
	var result *gorm.DB

	if result = db.Where("id = ? and user_id = ?", todoId, userId).Delete(&Todo{}); result.Error != nil {
		return common.CreateGenericInternalError(result.Error)
	}

	if result.RowsAffected == 0 {
		return common.CreateNotFoundError(fmt.Sprintf("Todo not found under given id (%d).", todoId))
	}

	return nil
}

func DeleteTodos(db *gorm.DB, userId int, todosId []int) *common.Error {
	var result *gorm.DB

	if result = db.Where("id in ? and user_id = ?", todosId, userId).Delete(&Todo{}); result.Error != nil {
		return common.CreateGenericInternalError(result.Error)
	}

	if result.RowsAffected == 0 {
		return common.CreateNotFoundError("Todos not found.")
	}

	return nil
}

func GetTodo(db *gorm.DB, userId int, todoId int) (*Todo, *common.Error) {
	var result *gorm.DB
	todo := Todo{}

	if result = db.Where("id = ? and user_id = ?", todoId, userId).First(&todo); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, common.CreateNotFoundError(fmt.Sprintf("Todo not found under given id (%d).", todoId))
		} else {
			return nil, common.CreateGenericInternalError(result.Error)
		}
	}

	return &todo, nil
}

func GetTodoChildren(db *gorm.DB, userId int, todoId int) ([]Todo, *common.Error) {
	var result *gorm.DB
	todos := []Todo{}

	if result = db.Where("parent_todo_id = ? and user_id = ?", todoId, userId).Order("created_at DESC").Find(&todos); result.Error != nil {
		return nil, common.CreateGenericInternalError(result.Error)
	}

	return todos, nil
}

func GetRootTodoChildren(db *gorm.DB, userId int) ([]Todo, *common.Error) {
	var result *gorm.DB
	todos := []Todo{}

	if result = db.Where("parent_todo_id is null and user_id = ?", userId).Order("created_at DESC").Find(&todos); result.Error != nil {
		return nil, common.CreateGenericInternalError(result.Error)
	}

	return todos, nil
}

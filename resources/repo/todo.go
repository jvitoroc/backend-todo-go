package repo

import (
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"
)

type Todo struct {
	ID           int64     `json:"todoId" db:"todoId"`
	ParentTodoID *int64    `json:"parentTodoId" db:"parentTodoId"`
	UserID       int64     `json:"userId" db:"userId"`
	Description  string    `json:"description" db:"description"`
	Completed    bool      `json:"completed" db:"completed"`
	CreatedAt    time.Time `json:"createdAt" db:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt" db:"updatedAt"`
}

func InsertTodo(db *sqlx.DB, userId int64, description string) (*int64, error) {
	now := time.Now()

	res, err := db.NamedExec("insert into todo (userId, description, createdAt, updatedAt) values (:userId, :description, :createdAt, :updatedAt)",
		map[string]interface{}{
			"userId":      userId,
			"description": description,
			"createdAt":   now,
			"updatedAt":   now})

	if err != nil {
		return nil, err
	}

	id, err := res.LastInsertId()

	return &id, err
}

func InsertTodoChild(db *sqlx.DB, userId int64, parentTodoId int64, description string) (*int64, error) {
	now := time.Now()

	res, err := db.NamedExec("insert into todo (userId, parentTodoId, description, createdAt, updatedAt) values (:userId, :parentTodoId, :description, :createdAt, :updatedAt)",
		map[string]interface{}{
			"userId":       userId,
			"parentTodoId": parentTodoId,
			"description":  description,
			"createdAt":    now,
			"updatedAt":    now})

	if err != nil {
		return nil, err
	}

	id, _ := res.LastInsertId()

	return &id, err
}

func UpdateTodo(db *sqlx.DB, userId int64, todoId int64, pairs map[string]interface{}) (sql.Result, error) {
	var command string = "update todo set "

	for key := range pairs {
		command = command + key + " = :" + key + ","
	}

	command = command[:len(command)-1]
	command = command + " where userId = :userId and todoId = :todoId"

	pairs["userId"] = userId
	pairs["todoId"] = todoId
	pairs["updatedAt"] = time.Now()

	res, err := db.NamedExec(command, pairs)

	return res, err
}

func DeleteTodo(db *sqlx.DB, userId int64, todoId int64) (sql.Result, error) {
	res, err := db.Exec("delete from todo where userId = $1 and todoId = $2", userId, todoId)

	return res, err
}

func DeleteTodos(db *sqlx.DB, userId int64, todosId []int64) (sql.Result, error) {
	query, args, err := sqlx.In("delete from todo where todoId in (?) and userId = ?", todosId, userId)
	query = db.Rebind(query)
	res, err := db.Exec(query, args...)

	return res, err
}

func GetTodo(db *sqlx.DB, userId int64, todoId int64) (*Todo, error) {
	todo := Todo{}

	err := db.Get(&todo, `select * from todo where userId = $1 and todoId = $2`, userId, todoId)

	return &todo, err
}

func GetTodoChildren(db *sqlx.DB, userId int64, todoId int64) ([]Todo, error) {
	todos := []Todo{}

	err := db.Select(&todos, `select * from todo where userId = $1 and parentTodoId = $2 order by createdAt desc`, userId, todoId)

	return todos, err
}

func GetRootTodoChildren(db *sqlx.DB, userId int64) ([]Todo, error) {
	todos := []Todo{}

	err := db.Select(&todos, `select * from todo where userId = $1 and parentTodoId is null order by createdAt desc`, userId)

	return todos, err
}

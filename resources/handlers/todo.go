package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/jvitoroc/todo-api/resources/repo"

	"github.com/gorilla/mux"
)

type TodoRequestBody struct {
	Description *string `json:"description"`
	Completed   *bool   `json:"completed"`
}

func initializeTodo(r *mux.Router) {
	addTodoHandlers(r)
}

func addTodoHandlers(r *mux.Router) {
	sr := r.PathPrefix("/todo").Subrouter()
	sr.Use(authenticateRequest)
	sr.Use(checkActivationState)
	sr.Handle("/", appHandler(createTodoHandler)).Methods("POST")
	sr.Handle("/{id:[0-9]+}", appHandler(createTodoHandler)).Methods("POST")
	sr.Handle("/{id:[0-9]+}", appHandler(updateTodoHandler)).Methods("PATCH")
	sr.Handle("/", appHandler(deleteManyTodosHandler)).Methods("DELETE")
	sr.Handle("/{id:[0-9]+}", appHandler(deleteTodoHandler)).Methods("DELETE")
	sr.Handle("/", appHandler(getTodoHandler)).Methods("GET")
	sr.Handle("/{id:[0-9]*}", appHandler(getTodoHandler)).Methods("GET")
}

func createTodoHandler(w http.ResponseWriter, r *http.Request) *appError {
	requestBody := TodoRequestBody{}

	if err := extractTodo(&requestBody, r); err != nil {
		return err
	}

	if err := validateTodo(&requestBody); err != nil {
		return err
	}

	userId, _ := strconv.ParseInt(r.Context().Value("userId").(string), 10, 64)
	parentTodoId, ok := mux.Vars(r)["id"]

	var err error
	var id *int64
	if ok {
		parentTodoId, _ := strconv.ParseInt(parentTodoId, 10, 64)
		id, err = repo.InsertTodoChild(db, userId, parentTodoId, *requestBody.Description)
	} else {
		id, err = repo.InsertTodo(db, userId, *requestBody.Description)
	}

	if err != nil {
		return unknownAppError(err)
	}

	todo, err := repo.GetTodo(db, userId, *id)

	if err != nil {
		return unknownAppError(err)
	}

	respond(
		map[string]interface{}{
			"message": "Todo successfully created.",
			"data":    todo,
		},
		http.StatusCreated,
		w,
	)
	return nil
}

func updateTodoHandler(w http.ResponseWriter, r *http.Request) *appError {
	userId, _ := strconv.ParseInt(r.Context().Value("userId").(string), 10, 64)
	todoId, _ := mux.Vars(r)["id"]
	_todoId, _ := strconv.ParseInt(todoId, 10, 64)

	requestBody := TodoRequestBody{}

	if appError := extractTodo(&requestBody, r); appError != nil {
		return appError
	}

	columns := make(map[string]interface{})

	if requestBody.Completed != nil {
		columns["completed"] = requestBody.Completed
	}

	if requestBody.Description != nil {
		columns["description"] = requestBody.Description
	}

	res, err := repo.UpdateTodo(db, userId, _todoId, columns)

	if rows, _err := res.RowsAffected(); _err == nil && rows == 0 {
		return createAppError(fmt.Sprintf(MSG_NOT_FOUND_ERROR, "Todo", _todoId), http.StatusNotFound)
	}

	if err != nil {
		return unknownAppError(err)
	}

	todo, err := repo.GetTodo(db, userId, _todoId)

	if err != nil {
		return unknownAppError(err)
	}

	respond(
		map[string]interface{}{
			"message": "Todo successfully updated.",
			"data":    todo,
		},
		http.StatusOK,
		w,
	)
	return nil
}

func deleteTodoHandler(w http.ResponseWriter, r *http.Request) *appError {
	userId, _ := strconv.ParseInt(r.Context().Value("userId").(string), 10, 64)
	todoId, _ := mux.Vars(r)["id"]
	_todoId, _ := strconv.ParseInt(todoId, 10, 64)

	res, err := repo.DeleteTodo(db, userId, _todoId)

	if err != nil {
		return unknownAppError(err)
	}

	if rows, err := res.RowsAffected(); err == nil && rows == 0 {
		return createAppError(fmt.Sprintf(MSG_NOT_FOUND_ERROR, "Todo", _todoId), http.StatusNotFound)
	}

	respondWithMessage("Todo successfully deleted.", http.StatusOK, w)
	return nil
}

func deleteManyTodosHandler(w http.ResponseWriter, r *http.Request) *appError {
	userId, _ := strconv.ParseInt(r.Context().Value("userId").(string), 10, 64)
	requestBody := make(map[string][]int64)

	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		return unknownAppError(err)
	}

	ids, ok := requestBody["ids"]

	if ok == false {
		return createAppError("A list of todos' IDs was not given.", http.StatusBadRequest)
	}

	if _, err := repo.DeleteTodos(db, userId, ids); err != nil {
		return unknownAppError(err)
	}

	respondWithMessage("Todos successfully deleted.", http.StatusOK, w)
	return nil
}

func getTodoHandler(w http.ResponseWriter, r *http.Request) *appError {
	userId, _ := strconv.ParseInt(r.Context().Value("userId").(string), 10, 64)
	todoId, ok := mux.Vars(r)["id"]
	_todoId, _ := strconv.ParseInt(todoId, 10, 64)

	var err error
	var todo *repo.Todo
	var parent *repo.Todo
	var grandParent *repo.Todo

	if ok {
		if todo, err = repo.GetTodo(db, userId, _todoId); err != nil {
			if err == sql.ErrNoRows {
				return createAppError(fmt.Sprintf(MSG_NOT_FOUND_ERROR, "Todo", _todoId), http.StatusNotFound)
			} else {
				return unknownAppError(err)
			}
		}
		if todo != nil && todo.ParentTodoID != nil {
			parent, err = repo.GetTodo(db, userId, *todo.ParentTodoID)
		}
		if parent != nil && parent.ParentTodoID != nil {
			grandParent, err = repo.GetTodo(db, userId, *parent.ParentTodoID)
		}
	}

	var children []repo.Todo

	if ok {
		children, err = repo.GetTodoChildren(db, userId, _todoId)
	} else {
		children, err = repo.GetRootTodoChildren(db, userId)
	}

	if err != nil {
		return unknownAppError(err)
	}

	respond(
		map[string]interface{}{
			"message": "Todo successfully retrieved.",
			"data": map[string]interface{}{
				"grandParent": grandParent,
				"parent":      parent,
				"todo":        todo,
				"children":    children,
			},
		},
		http.StatusOK,
		w,
	)
	return nil
}

func extractTodo(todo *TodoRequestBody, r *http.Request) *appError {
	err := json.NewDecoder(r.Body).Decode(todo)
	if err != nil {
		return createAppError(fmt.Sprintf(MSG_UNKNOWN_ERROR, err.Error()), http.StatusBadRequest)
	}
	return nil
}

func validateTodo(todo *TodoRequestBody) *appError {
	errors := map[string]string{}

	if todo.Description == nil || *todo.Description == "" {
		errors["description"] = "Description field is empty or missing."
	}

	if len(errors) == 0 {
		return nil
	} else {
		return createMappedAppError(MSG_ONE_MORE_ERRORS, errors, http.StatusBadRequest)
	}
}

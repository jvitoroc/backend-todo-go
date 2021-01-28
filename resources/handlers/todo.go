package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/jvitoroc/todo-api/resources/common"
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

func createTodoHandler(w http.ResponseWriter, r *http.Request) *common.Error {
	requestBody := TodoRequestBody{}

	if err := extractTodo(&requestBody, r); err != nil {
		return err
	}

	if err := validateTodo(&requestBody); err != nil {
		return err
	}

	var err *common.Error
	var todo *repo.Todo

	userId, _ := strconv.Atoi(r.Context().Value("userId").(string))
	parentTodoId, ok := mux.Vars(r)["id"]

	if ok { // check whether the user provided a parent to the new todo
		parentTodoId, _ := strconv.Atoi(parentTodoId)
		todo, err = repo.InsertTodoChild(db, userId, parentTodoId, *requestBody.Description)
	} else {
		todo, err = repo.InsertTodo(db, userId, *requestBody.Description)
	}

	if err != nil {
		return err
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

func updateTodoHandler(w http.ResponseWriter, r *http.Request) *common.Error {
	userId, _ := extractContextInt("userId", r)
	todoId, _ := extractParamInt("id", r)

	requestBody := TodoRequestBody{}

	if appError := extractTodo(&requestBody, r); appError != nil {
		return appError
	}

	columns := make(map[string]interface{})

	if requestBody.Completed != nil {
		columns["completed"] = *requestBody.Completed
	}

	if requestBody.Description != nil {
		columns["description"] = *requestBody.Description
	}

	if err := repo.UpdateTodo(db, userId, todoId, columns); err != nil {
		return err
	}

	var todo *repo.Todo
	var err *common.Error

	if todo, err = repo.GetTodo(db, userId, todoId); err != nil {
		return err
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

func deleteTodoHandler(w http.ResponseWriter, r *http.Request) *common.Error {
	userId, _ := extractContextInt("userId", r)
	todoId, _ := extractParamInt("id", r)

	if err := repo.DeleteTodo(db, userId, todoId); err != nil {
		return err
	}

	respondWithMessage("Todo successfully deleted.", http.StatusOK, w)
	return nil
}

func deleteManyTodosHandler(w http.ResponseWriter, r *http.Request) *common.Error {
	userId, _ := extractContextInt("userId", r)
	requestBody := make(map[string][]int)

	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		return common.CreateGenericBadRequestError(err)
	}

	ids, ok := requestBody["ids"]

	if ok == false {
		return common.CreateBadRequestError("A list of todos' IDs was not given.")
	}

	if err := repo.DeleteTodos(db, userId, ids); err != nil {
		return err
	}

	respondWithMessage("Todos successfully deleted.", http.StatusOK, w)
	return nil
}

func getTodoHandler(w http.ResponseWriter, r *http.Request) *common.Error {
	userId, _ := extractContextInt("userId", r)
	todoId, ok := extractParamInt("id", r)

	var err *common.Error
	var todo *repo.Todo
	var parent *repo.Todo
	var grandParent *repo.Todo

	if ok {
		if todo, err = repo.GetTodo(db, userId, todoId); err != nil {
			return err
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
		children, err = repo.GetTodoChildren(db, userId, todoId)
	} else {
		children, err = repo.GetRootTodoChildren(db, userId)
	}

	if err != nil {
		return err
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

func extractTodo(todo *TodoRequestBody, r *http.Request) *common.Error {
	err := json.NewDecoder(r.Body).Decode(todo)
	if err != nil {
		return common.CreateGenericBadRequestError(err)
	}
	return nil
}

func validateTodo(todo *TodoRequestBody) *common.Error {
	errors := map[string]string{}

	if todo.Description == nil || *todo.Description == "" {
		errors["description"] = "Description field is empty or missing."
	}

	if len(errors) == 0 {
		return nil
	} else {
		return common.CreateFormError(errors)
	}
}

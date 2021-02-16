package api

import (
	"net/http"

	hn "github.com/jvitoroc/todo-go/api/handler"
	"github.com/jvitoroc/todo-go/model"
	"github.com/jvitoroc/todo-go/util"
)

func (api *API) InitTodo() {
	api.Router.Todo.Handle("/{todoId:[0-9]*}", api.createProtectedHandler(api.CreateTodo, true)).Methods("POST")
	api.Router.Todo.Handle("", api.createProtectedHandler(api.GetRootTodoChildren, true)).Methods("GET")
	api.Router.Todo.Handle("/{todoId:[0-9]+}", api.createProtectedHandler(api.GetTodo, true)).Methods("GET")
	api.Router.Todo.Handle("/{todoId:[0-9]+}", api.createProtectedHandler(api.UpdateTodo, true)).Methods("PATCH")
	api.Router.Todo.Handle("/{todoId:[0-9]+}", api.createProtectedHandler(api.DeleteTodo, true)).Methods("DELETE")
	api.Router.Todo.Handle("", api.createProtectedHandler(api.DeleteManyTodos, true)).Methods("DELETE")
}

func (api *API) CreateTodo(hn *hn.RequestContext, w http.ResponseWriter, r *http.Request) hn.Response {
	todo, err := model.TodoFromJson(r.Body)
	if err != nil {
		return err
	}

	if err := todo.Validate(); err != nil {
		return err
	}

	todoId, ok := util.ExtractParamInt("todoId", r)
	if ok {
		todo.ParentTodoID = &todoId
	}

	todo.UserID = hn.CurrentUser.ID

	if err := api.App.CreateTodo(todo); err != nil {
		return err
	}

	return model.NewCreatedResponse(model.MSG_TODO_CREATED).AddObject("todo", todo)
}

func (api *API) GetTodo(hn *hn.RequestContext, w http.ResponseWriter, r *http.Request) hn.Response {
	todoId, _ := util.ExtractParamInt("todoId", r)
	todo, err := api.App.GetTodo(todoId, hn.CurrentUser.ID)
	if err != nil {
		return err
	}

	parentsCount, _ := util.ExtractFormInt("parents-count", r)
	if parentsCount < 0 {
		parentsCount = 0
	}

	parents := make([]model.Todo, 0) // allocate array so we can return at least a empty array to the client
	if todo.ParentTodoID != nil {
		parents, err = api.App.GetTodosBottomToTop(*todo.ParentTodoID, hn.CurrentUser.ID, parentsCount)
		if err != nil {
			return err
		}
	}

	todos, err := api.App.GetTodoChildren(todoId, hn.CurrentUser.ID)
	if err != nil {
		return err
	}

	res := model.NewOKResponse(model.MSG_TODO_RETRIEVED)
	res.AddObject("children", todos)
	res.AddObject("todo", todo)
	res.AddObject("parents", parents)

	return res
}

func (api *API) GetRootTodoChildren(ctx *hn.RequestContext, w http.ResponseWriter, r *http.Request) hn.Response {
	todos, err := api.App.GetRootTodoChildren(ctx.CurrentUser.ID)
	if err != nil {
		return err
	}

	return model.NewOKResponse(model.MSG_TODO_RETRIEVED).AddObject("children", todos)
}

func (api *API) UpdateTodo(ctx *hn.RequestContext, w http.ResponseWriter, r *http.Request) hn.Response {
	todo, err := model.UpdateTodoFromJson(r.Body)
	if err != nil {
		return err
	}

	if err := todo.Validate(); err != nil {
		return err
	}

	todoId, _ := util.ExtractParamInt("todoId", r)
	todo.ID = todoId
	dbTodo, err := api.App.UpdateTodo(todo, ctx.CurrentUser.ID)
	if err != nil {
		return err
	}

	return model.NewOKResponse(model.MSG_TODO_UPDATED).AddObject("todo", dbTodo)
}

func (api *API) DeleteTodo(ctx *hn.RequestContext, w http.ResponseWriter, r *http.Request) hn.Response {
	todoId, _ := util.ExtractParamInt("todoId", r)
	err := api.App.DeleteTodo(todoId, ctx.CurrentUser.ID)
	if err != nil {
		return err
	}

	return model.NewOKResponse(model.MSG_TODO_DELETED)
}

func (api *API) DeleteManyTodos(ctx *hn.RequestContext, w http.ResponseWriter, r *http.Request) hn.Response {
	delete, err := model.DeleteManyTodosFromJson(r.Body)
	if err != nil {
		return err
	}

	if delete.Validate(); err != nil {
		return err
	}

	if err := api.App.DeleteManyTodos(delete.IDs, ctx.CurrentUser.ID); err != nil {
		return err
	}

	return model.NewOKResponse(model.MSG_TODOS_DELETED)
}

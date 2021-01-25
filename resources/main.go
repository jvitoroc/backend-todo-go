package resources

import (
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"github.com/jvitoroc/todo-api/resources/handlers"
)

var schema = `
DROP TABLE IF EXISTS todo;
DROP TABLE IF EXISTS user_activation_request;
DROP TABLE IF EXISTS user;

CREATE TABLE user(
	userId integer primary key,
	username text unique,
	email text unique,
	passwordHash text,
	active bool default false,
	createdAt datetime,
	updatedAt datetime
);

CREATE TABLE user_activation_request(
	userId integer primary key,
	expiresAt datetime,
	code text,
	FOREIGN KEY(userId) REFERENCES user(userId) ON DELETE CASCADE
);

CREATE TABLE todo(
	todoId integer primary key,
	parentTodoId integer null,
	userId integer,
	description text,
	completed bool default false,
	createdAt datetime,
	updatedAt datetime,
	FOREIGN KEY(userId) REFERENCES user(userId) ON DELETE CASCADE,
	FOREIGN KEY(parentTodoId) REFERENCES todo(todoId) ON DELETE CASCADE
);
`

func Initialize(r *mux.Router, db *sqlx.DB) {
	// db.MustExec(schema)
	handlers.Initialize(r, db)
}

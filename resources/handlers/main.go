package handlers

import (
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
)

var db *sqlx.DB

func Initialize(r *mux.Router, _db *sqlx.DB) {
	db = _db
	initializeUser(r)
	initializeTodo(r)
}

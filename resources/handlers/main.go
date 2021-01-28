package handlers

import (
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

var db *gorm.DB

func Initialize(r *mux.Router, _db *gorm.DB) {
	db = _db
	initializeUser(r)
	initializeTodo(r)
}

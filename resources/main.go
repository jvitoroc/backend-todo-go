package resources

import (
	"github.com/gorilla/mux"
	"github.com/jvitoroc/todo-api/resources/handlers"
	"github.com/jvitoroc/todo-api/resources/repo"
	"gorm.io/gorm"
)

func Initialize(r *mux.Router, db *gorm.DB) {
	repo.Initialize(db)
	handlers.Initialize(r, db)
}

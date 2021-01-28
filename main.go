package main

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/jvitoroc/todo-api/resources"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setBasics(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func corsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, PATCH, DELETE")
}

func main() {
	err := godotenv.Load()
	db, err := gorm.Open(sqlite.Open("file:test.db?&cache=shared&_fk=1"), &gorm.Config{})

	if err != nil {
		panic(err)
	}

	r := mux.NewRouter()
	r.StrictSlash(true)
	r.PathPrefix("/").HandlerFunc(corsHandler).Methods("OPTIONS")
	r.Use(setBasics)

	// init routes and models
	resources.Initialize(r, db)

	http.ListenAndServe(":8000", r)
}

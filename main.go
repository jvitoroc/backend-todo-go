package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/jvitoroc/todo-api/resources"
	"google.golang.org/api/idtoken"
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

	payload, err := idtoken.Validate(
		context.Background(),
		"eyJhbGciOiJSUzI1NiIsImtpZCI6IjAzYjJkMjJjMmZlY2Y4NzNlZDE5ZTViOGNmNzA0YWZiN2UyZWQ0YmUiLCJ0eXAiOiJKV1QifQ.eyJpc3MiOiJhY2NvdW50cy5nb29nbGUuY29tIiwiYXpwIjoiMTAxMTIwMDg0ODYyOS1uMnMzMHIzNTNrN2h2Z2s5azczYnNqMjBiOTExb212di5hcHBzLmdvb2dsZXVzZXJjb250ZW50LmNvbSIsImF1ZCI6IjEwMTEyMDA4NDg2MjktbjJzMzByMzUzazdodmdrOWs3M2JzajIwYjkxMW9tdnYuYXBwcy5nb29nbGV1c2VyY29udGVudC5jb20iLCJzdWIiOiIxMDM0NTQ3MDY5NjAyNDQ2NjE1NDEiLCJlbWFpbCI6Imp2aXRvcm9jMTdAZ21haWwuY29tIiwiZW1haWxfdmVyaWZpZWQiOnRydWUsImF0X2hhc2giOiJaZEZTNFhyMEhzMkRvY3JSSWxjU3JnIiwibmFtZSI6Ikpvw6NvIFZpdG9yIGRlIE9saXZlaXJhIENhcmxvcyIsInBpY3R1cmUiOiJodHRwczovL2xoMy5nb29nbGV1c2VyY29udGVudC5jb20vLThibTU0RXh6eGgwL0FBQUFBQUFBQUFJL0FBQUFBQUFBQUFBL0FNWnV1Y21rTnR5OWEwWFdiNWI5TlJYRHFOd01FbUc0ZkEvczk2LWMvcGhvdG8uanBnIiwiZ2l2ZW5fbmFtZSI6Ikpvw6NvIFZpdG9yIiwiZmFtaWx5X25hbWUiOiJkZSBPbGl2ZWlyYSBDYXJsb3MiLCJsb2NhbGUiOiJwdC1CUiIsImlhdCI6MTYxMTg3ODY1OCwiZXhwIjoxNjExODgyMjU4LCJqdGkiOiJkMjQzNDkzYzQ4OGYwNWQyYTczOTcxNTViMzM0YWQwNTg0YWI2ZDEwIn0.mXaTEzDP3DU-5atWdWZzPM-zbX-Y22fFpt6A917bR0gbXcvq-GQI7ksUcpzu1tuKsvt9aecwiEMvdS9QI5iOKg9tFxjTbqh_zuoxtob4O-d7WbxtmvHlbAItTvfTFCvZE7Qii_IscaNpksG6pj3ETrbHBm9wrY6yU7OQnyjF1bnubx-PcSXAAcebaHnWq6KpxuGHXRayhocMgQ7cVHcwdrjUFm-iJTsSkbTwdXjLAxHtqJtLUHTEoxAaQ7shFUc0JVWruElFNBde66uAoM8B8yDf6pXvuz9pu6IiHHn1UAkLRzzekPjjbnCUrt5urnNFLjHtZEc0gy0WFkvAbmtp3Q",
		os.Getenv("GOOGLE_CLIENT_ID"),
	)

	if err != nil {
		// Not a Google ID token.
	}
	log.Print(payload.Claims["sub"])

	// init routes and models
	resources.Initialize(r, db)

	http.ListenAndServe(":8000", r)
}

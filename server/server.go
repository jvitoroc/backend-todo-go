package server

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jvitoroc/todo-go/api"
	"github.com/jvitoroc/todo-go/app"
	"github.com/jvitoroc/todo-go/auth"
	"github.com/jvitoroc/todo-go/config"
	"github.com/jvitoroc/todo-go/email"
	"github.com/jvitoroc/todo-go/repository"
	"github.com/rs/cors"
)

type Server struct {
	API    *api.API
	Config *config.Config
	Router *mux.Router
}

func NewServer() *Server {
	cfg := config.NewConfig(config.CFG_PROD)
	repo := repository.NewRepository(cfg)
	email := email.NewEmailService(cfg)
	auth := auth.NewAuthService(cfg)
	router := mux.NewRouter()

	app := app.NewApp(repo, email, auth, cfg)
	api := api.NewAPI(app, router)

	return &Server{
		API:    api,
		Config: cfg,
		Router: router,
	}
}

func (s *Server) Start() {
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedHeaders: []string{"*"},
		AllowedMethods: []string{"OPTIONS", "POST", "PUT", "GET", "PATCH", "DELETE"},
	})

	addr := ":" + s.Config.Server.Port
	s.Router.Use(setBasicsMiddleware)

	if err := http.ListenAndServe(addr, c.Handler(s.Router)); err != nil {
		log.Fatalf("Could not start the server: %s", err.Error())
	}
}

func setBasicsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

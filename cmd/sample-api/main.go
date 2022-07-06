package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jmoiron/sqlx"
	"github.com/podhmo-sandbox/sample-api/pkg/dblib"
	"github.com/podhmo-sandbox/sample-api/repository"
	"github.com/podhmo-sandbox/sample-api/webapi/todo"
	"github.com/podhmo/flagstruct"
	_ "modernc.org/sqlite"
)

type Config struct {
	DB   dblib.Config `flag:"db"`
	Addr string       `flag:"addr"`
}

func mount(r chi.Router, db *sqlx.DB) {
	r.Route("/todos", func(r chi.Router) {
		repo := repository.NewTodoRepository(db)
		r.MethodFunc("GET", "/", todo.GetTodos(repo))
		r.MethodFunc("POST", "/", todo.PostTodo(repo))
		r.MethodFunc("PUT", "/", todo.PutTodo(repo))
		r.MethodFunc("DELETE", "/", todo.DeleteTodo(repo))
	})
}

func run(config Config) error {
	server := http.Server{
		Addr: config.Addr,
	}

	r := chi.NewRouter()

	// A good base middleware stack
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	json.NewEncoder(os.Stdout).Encode(config)
	db := sqlx.MustConnect(config.DB.DSN, config.DB.Driver)
	mount(r, db)
	return server.ListenAndServe()
}

func main() {
	config := &Config{DB: dblib.DefaultConfig(), Addr: ":8080"} // default values
	flagstruct.Parse(config)
	if err := run(*config); err != nil {
		log.Fatalf("!! %+v", err)
	}
}

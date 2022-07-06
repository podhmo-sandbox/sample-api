package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jmoiron/sqlx"
	"github.com/podhmo-sandbox/sample-api/repository"
	"github.com/podhmo-sandbox/sample-api/webapi/todo"
	_ "modernc.org/sqlite"
)

func mount(r chi.Router, db *sqlx.DB) {
	r.Route("/todos", func(r chi.Router) {
		repo := repository.NewTodoRepository(db)
		r.MethodFunc("GET", "", todo.GetTodos(repo))
		r.MethodFunc("POST", "", todo.PostTodo(repo))
		r.MethodFunc("PUT", "", todo.PutTodo(repo))
		r.MethodFunc("DELETE", "", todo.DeleteTodo(repo))
	})
}

func main() {
	server := http.Server{
		Addr: ":8080",
	}

	r := chi.NewRouter()

	// A good base middleware stack
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	db := sqlx.MustConnect("sqlite", ":memory:")
	mount(r, db)
	server.ListenAndServe()
}

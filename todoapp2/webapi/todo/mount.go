package todo

import (
	"github.com/go-chi/chi/v5"
	"github.com/podhmo/quickapi"
)

func Mount(r chi.Router, repo todoRepository) chi.Router {
	r.Route("/todos", func(r chi.Router) {
		h := &Handler{Repo: repo}
		r.MethodFunc("GET", "/", quickapi.Lift(h.GetTodos))
		r.MethodFunc("POST", "/", quickapi.Lift(h.PostTodo))
		r.MethodFunc("PUT", "/{todoId}", quickapi.Lift(h.PutTodo))
		r.MethodFunc("DELETE", "/{todoId}", quickapi.Lift(h.DeleteTodo))
	})
	return r
}

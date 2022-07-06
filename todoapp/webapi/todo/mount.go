package todo

import (
	"github.com/go-chi/chi/v5"
)

func Mount(r chi.Router, repo todoRepository) {
	r.Route("/todos", func(r chi.Router) {
		r.MethodFunc("GET", "/", GetTodos(repo))
		r.MethodFunc("POST", "/", PostTodo(repo))
		r.MethodFunc("PUT", "/{todoId}", PutTodo(repo))
		r.MethodFunc("DELETE", "/{todoId}", DeleteTodo(repo))
	})
}

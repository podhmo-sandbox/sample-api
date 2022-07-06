package controller

import (
	"github.com/go-chi/chi/v5"
)

func Mount(r chi.Router, c *TodoController) {
	r.Route("/todos", func(r chi.Router) {
		r.MethodFunc("GET", "", c.GetTodos)
		r.MethodFunc("POST", "", c.PostTodo)
		r.MethodFunc("PUT", "", c.PutTodo)
		r.MethodFunc("DELETE", "", c.DeleteTodo)
	})
}

package todo

import (
	"github.com/podhmo/quickapi/experimental/define"
)

func Mount(
	repo todoRepository,
) func(*define.BuildContext) {
	return func(bc *define.BuildContext) {
		h := &Handler{Repo: repo}
		define.Get(bc, "/todos", h.GetTodos)
		define.Post(bc, "/todos", h.PostTodo)
		define.Put(bc, "/todos/{todoId}", h.PutTodo)
		define.Delete(bc, "/todos/{todoId}", h.DeleteTodo)
	}
}

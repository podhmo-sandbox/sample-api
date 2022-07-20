package todo

import (
	"context"
	"strconv"

	"github.com/podhmo-sandbox/sample-api/todoapp2/entity"
	"github.com/podhmo/quickapi"
)

type todoRepository interface {
	GetTodos() (todos []entity.Todo, err error)
	InsertTodo(todo entity.Todo) (id int, err error)
	UpdateTodo(todo entity.Todo) (err error)
	DeleteTodo(id int) (err error)
}

type Handler struct {
	Repo todoRepository
}

func (h *Handler) GetTodos(context.Context, quickapi.Empty) (output TodosResponse, err error) {
	todos, err := h.Repo.GetTodos()
	if err != nil {
		return output, err
	}

	for _, v := range todos {
		output.Todos = append(output.Todos, TodoResponse{ID: v.ID, Title: v.Title, Content: v.Content})
	}
	return output, nil
}

func (h *Handler) PostTodo(ctx context.Context, input TodoRequest) (output any, err error) {
	todo := entity.Todo{Title: input.Title, Content: input.Content}
	id, err := h.Repo.InsertTodo(todo)
	if err != nil {
		return output, err
	}

	r := quickapi.GetRequest(ctx)
	return output, quickapi.Redirect(201, r.Host+r.URL.Path+strconv.Itoa(id))
}

type putTodoInput struct {
	TodoID int `path:"todoId" openapi:"path"`
	TodoRequest
}

func (h *Handler) PutTodo(ctx context.Context, input putTodoInput) (output any, err error) {
	todo := entity.Todo{ID: input.TodoID, Title: input.Title, Content: input.Content}
	if err = h.Repo.UpdateTodo(todo); err != nil {
		return output, err
	}
	return quickapi.NoContent(204), nil
}

type deleteTodoInput struct {
	TodoID int `path:"todoId" openapi:"path"`
}

func (h *Handler) DeleteTodo(ctx context.Context, input deleteTodoInput) (output any, err error) {
	if err = h.Repo.DeleteTodo(input.TodoID); err != nil {
		return output, err
	}
	return quickapi.NoContent(204), nil
}

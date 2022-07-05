package controller

import (
	"encoding/json"
	"net/http"
	"path"
	"strconv"

	"github.com/podhmo-sandbox/sample-api/controller/dto"
	"github.com/podhmo-sandbox/sample-api/model/entity"
)

type TodoRepository interface {
	GetTodos() (todos []entity.Todo, err error)
	InsertTodo(todo entity.Todo) (id int, err error)
	UpdateTodo(todo entity.Todo) (err error)
	DeleteTodo(id int) (err error)
}

type TodoController struct {
	tr TodoRepository
}

func NewTodoController(tr TodoRepository) *TodoController {
	return &TodoController{tr}
}

func (tc *TodoController) GetTodos(w http.ResponseWriter, r *http.Request) {
	todos, err := tc.tr.GetTodos()
	if err != nil {
		w.WriteHeader(500)
		return
	}

	var todoResponses []dto.TodoResponse
	for _, v := range todos {
		todoResponses = append(todoResponses, dto.TodoResponse{ID: v.ID, Title: v.Title, Content: v.Content})
	}

	var todosResponse dto.TodosResponse
	todosResponse.Todos = todoResponses

	output, _ := json.MarshalIndent(todosResponse.Todos, "", "\t\t")

	w.Header().Set("Content-Type", "application/json")
	w.Write(output)
}

func (tc *TodoController) PostTodo(w http.ResponseWriter, r *http.Request) {
	body := make([]byte, r.ContentLength)
	r.Body.Read(body)
	var todoRequest dto.TodoRequest
	json.Unmarshal(body, &todoRequest)

	todo := entity.Todo{Title: todoRequest.Title, Content: todoRequest.Content}
	id, err := tc.tr.InsertTodo(todo)
	if err != nil {
		w.WriteHeader(500)
		return
	}

	w.Header().Set("Location", r.Host+r.URL.Path+strconv.Itoa(id))
	w.WriteHeader(201)
}

func (tc *TodoController) PutTodo(w http.ResponseWriter, r *http.Request) {
	todoID, err := strconv.Atoi(path.Base(r.URL.Path))
	if err != nil {
		w.WriteHeader(400)
		return
	}

	body := make([]byte, r.ContentLength)
	r.Body.Read(body)
	var todoRequest dto.TodoRequest
	json.Unmarshal(body, &todoRequest)

	todo := entity.Todo{ID: todoID, Title: todoRequest.Title, Content: todoRequest.Content}
	err = tc.tr.UpdateTodo(todo)
	if err != nil {
		w.WriteHeader(500)
		return
	}

	w.WriteHeader(204)
}

func (tc *TodoController) DeleteTodo(w http.ResponseWriter, r *http.Request) {
	todoID, err := strconv.Atoi(path.Base(r.URL.Path))
	if err != nil {
		w.WriteHeader(400)
		return
	}

	err = tc.tr.DeleteTodo(todoID)
	if err != nil {
		w.WriteHeader(500)
		return
	}

	w.WriteHeader(204)
}

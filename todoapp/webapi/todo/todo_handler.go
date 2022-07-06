package todo

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/podhmo-sandbox/sample-api/todoapp/entity"
)

type todoRepository interface {
	GetTodos() (todos []entity.Todo, err error)
	InsertTodo(todo entity.Todo) (id int, err error)
	UpdateTodo(todo entity.Todo) (err error)
	DeleteTodo(id int) (err error)
}

func GetTodos(repo todoRepository) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		todos, err := repo.GetTodos()
		if err != nil {
			w.WriteHeader(500)
			return
		}

		var todoResponses []TodoResponse
		for _, v := range todos {
			todoResponses = append(todoResponses, TodoResponse{ID: v.ID, Title: v.Title, Content: v.Content})
		}

		var todosResponse TodosResponse
		todosResponse.Todos = todoResponses

		output, _ := json.MarshalIndent(todosResponse, "", "\t\t")

		w.Header().Set("Content-Type", "application/json")
		w.Write(output)
	}
}

func PostTodo(repo todoRepository) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		body := make([]byte, r.ContentLength)
		r.Body.Read(body)
		var todoRequest TodoRequest
		json.Unmarshal(body, &todoRequest)

		todo := entity.Todo{Title: todoRequest.Title, Content: todoRequest.Content}
		id, err := repo.InsertTodo(todo)
		if err != nil {
			w.WriteHeader(500)
			return
		}

		w.Header().Set("Location", r.Host+r.URL.Path+strconv.Itoa(id))
		w.WriteHeader(201)
	}
}

func PutTodo(repo todoRepository) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		todoID, err := strconv.Atoi(chi.URLParam(r, "todoId"))
		if err != nil {
			w.WriteHeader(400)
			return
		}

		body := make([]byte, r.ContentLength)
		r.Body.Read(body)
		var todoRequest TodoRequest
		json.Unmarshal(body, &todoRequest)

		todo := entity.Todo{ID: todoID, Title: todoRequest.Title, Content: todoRequest.Content}
		err = repo.UpdateTodo(todo)
		if err != nil {
			w.WriteHeader(500)
			return
		}

		w.WriteHeader(204)
	}
}

func DeleteTodo(repo todoRepository) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		todoID, err := strconv.Atoi(chi.URLParam(r, "todoId"))
		if err != nil {
			w.WriteHeader(400)
			return
		}

		err = repo.DeleteTodo(todoID)
		if err != nil {
			w.WriteHeader(500)
			return
		}

		w.WriteHeader(204)
	}
}

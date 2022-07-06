package controller

import (
	"encoding/json"
	"net/http"
	"path"
	"strconv"

	"github.com/podhmo-sandbox/sample-api/controller/dto"
	"github.com/podhmo-sandbox/sample-api/model/entity"
)

func GetTodos(repo interface {
	GetTodos() (todos []entity.Todo, err error)
}) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		todos, err := repo.GetTodos()
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
}

func PostTodo(repo interface {
	InsertTodo(todo entity.Todo) (id int, err error)
}) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		body := make([]byte, r.ContentLength)
		r.Body.Read(body)
		var todoRequest dto.TodoRequest
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

func PutTodo(repo interface {
	UpdateTodo(todo entity.Todo) (err error)
}) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
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
		err = repo.UpdateTodo(todo)
		if err != nil {
			w.WriteHeader(500)
			return
		}

		w.WriteHeader(204)
	}
}

func DeleteTodo(repo interface {
	DeleteTodo(id int) (err error)
}) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		todoID, err := strconv.Atoi(path.Base(r.URL.Path))
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

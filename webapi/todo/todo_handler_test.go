package todo_test

import (
	"encoding/json"
	"net/http/httptest"
	"strings"
	"testing"

	webapi "github.com/podhmo-sandbox/sample-api/webapi/todo"
)

func TestGetTodos_NotFound(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/todos/", nil)

	target := webapi.GetTodos(&MockTodoRepository{})
	target(w, r)

	if w.Code != 200 {
		t.Errorf("Response cod is %v", w.Code)
	}
	if w.Header().Get("Content-Type") != "application/json" {
		t.Errorf("Content-Type is %v", w.Header().Get("Content-Type"))
	}

	body := make([]byte, w.Body.Len())
	w.Body.Read(body)
	var todosResponse webapi.TodosResponse
	json.Unmarshal(body, &todosResponse)
	if len(todosResponse.Todos) != 0 {
		t.Errorf("Response is %v", todosResponse.Todos)
	}
}

func TestGetTodos_ExistTodo(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/todos/", nil)

	target := webapi.GetTodos(&MockTodoRepositoryGetTodosExist{})
	target(w, r)

	if w.Code != 200 {
		t.Errorf("Response cod is %v", w.Code)
	}
	if w.Header().Get("Content-Type") != "application/json" {
		t.Errorf("Content-Type is %v", w.Header().Get("Content-Type"))
	}

	body := make([]byte, w.Body.Len())
	w.Body.Read(body)
	var todosResponse webapi.TodosResponse
	json.Unmarshal(body, &todosResponse.Todos)
	if len(todosResponse.Todos) != 2 {
		t.Errorf("Response is %v", todosResponse.Todos)
	}
}

func TestGetTodos_Error(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/todos/", nil)

	target := webapi.GetTodos(&MockTodoRepositoryError{})
	target(w, r)

	if w.Code != 500 {
		t.Errorf("Response cod is %v", w.Code)
	}
	if w.Header().Get("Content-Type") != "" {
		t.Errorf("Content-Type is %v", w.Header().Get("Content-Type"))
	}

	if w.Body.Len() != 0 {
		t.Errorf("body is %v", w.Body.Len())
	}
}

func TestPostTodo_OK(t *testing.T) {
	json := strings.NewReader(`{"title":"test-title","content":"test-content"}`)
	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/todos/", json)

	target := webapi.PostTodo(&MockTodoRepository{})
	target(w, r)

	if w.Code != 201 {
		t.Errorf("Response cod is %v", w.Code)
	}
	if w.Header().Get("Location") != r.Host+r.URL.Path+"2" {
		t.Errorf("Location is %v", w.Header().Get("Location"))
	}
}

func TestPostTodo_Error(t *testing.T) {
	json := strings.NewReader(`{"title":"test-title","contents":"test-content"}`)
	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/todos/", json)

	target := webapi.PostTodo(&MockTodoRepositoryError{})
	target(w, r)

	if w.Code != 500 {
		t.Errorf("Response cod is %v", w.Code)
	}
	if w.Header().Get("Location") != "" {
		t.Errorf("Location is %v", w.Header().Get("Location"))
	}
}

func TestPutTodo_OK(t *testing.T) {
	json := strings.NewReader(`{"title":"test-title","contents":"test-content"}`)
	w := httptest.NewRecorder()
	r := httptest.NewRequest("PUT", "/todos/2", json)

	target := webapi.PutTodo(&MockTodoRepository{})
	target(w, r)

	if w.Code != 204 {
		t.Errorf("Response cod is %v", w.Code)
	}
}

func TestPutTodo_InvalidPath(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest("PUT", "/todos/", nil)

	target := webapi.PutTodo(&MockTodoRepository{})
	target(w, r)

	if w.Code != 400 {
		t.Errorf("Response cod is %v", w.Code)
	}
}

func TestPutTodo_Error(t *testing.T) {
	json := strings.NewReader(`{"title":"test-title","contents":"test-content"}`)
	w := httptest.NewRecorder()
	r := httptest.NewRequest("PUT", "/todos/2", json)

	target := webapi.PutTodo(&MockTodoRepositoryError{})
	target(w, r)

	if w.Code != 500 {
		t.Errorf("Response cod is %v", w.Code)
	}
}

func TestDeleteTodo_OK(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest("DELETE", "/todos/2", nil)

	target := webapi.DeleteTodo(&MockTodoRepository{})
	target(w, r)

	if w.Code != 204 {
		t.Errorf("Response cod is %v", w.Code)
	}
}

func TestDeleteTodo_InvalidPath(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest("DELETE", "/todos/", nil)

	target := webapi.DeleteTodo(&MockTodoRepositoryError{})
	target(w, r)

	if w.Code != 400 {
		t.Errorf("Response cod is %v", w.Code)
	}
}

func TestDeleteTodo_Error(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest("DELETE", "/todos/2", nil)

	target := webapi.DeleteTodo(&MockTodoRepositoryError{})
	target(w, r)

	if w.Code != 500 {
		t.Errorf("Response cod is %v", w.Code)
	}
}

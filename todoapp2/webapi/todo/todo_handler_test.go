package todo_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/google/go-cmp/cmp"
	"github.com/podhmo-sandbox/sample-api/todoapp2/webapi/todo"
	webapi "github.com/podhmo-sandbox/sample-api/todoapp2/webapi/todo"
	"github.com/podhmo/quickapi"
	"github.com/podhmo/quickapi/quickapitest"
)

// TODO: performance up

func TestGetTodos(t *testing.T) {
	t.Run("not-found", func(t *testing.T) {
		handler := http.HandlerFunc(webapi.Mount(chi.NewRouter(), &MockTodoRepository{}).ServeHTTP)
		code := 200
		want := webapi.TodosResponse{Todos: []todo.TodoResponse{}}

		req := httptest.NewRequest("GET", "/todos/", nil)
		got := quickapitest.DoRequest[webapi.TodosResponse](t, req, code, handler)

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("%s %s, response mismatch (-want +got):\n%s", req.Method, req.URL.Path, diff)
		}
	})

	t.Run("ok", func(t *testing.T) {
		handler := http.HandlerFunc(webapi.Mount(chi.NewRouter(), &MockTodoRepositoryGetTodosExist{}).ServeHTTP)
		code := 200
		want := webapi.TodosResponse{Todos: []todo.TodoResponse{
			{ID: 1, Title: "title1", Content: "contents1"},
			{ID: 2, Title: "title2", Content: "contents2"},
		}}

		req := httptest.NewRequest("GET", "/todos/", nil)
		got := quickapitest.DoRequest[webapi.TodosResponse](t, req, code, handler)

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("%s %s, response mismatch (-want +got):\n%s", req.Method, req.URL.Path, diff)
		}
	})

	t.Run("error", func(t *testing.T) {
		handler := http.HandlerFunc(webapi.Mount(chi.NewRouter(), &MockTodoRepositoryError{}).ServeHTTP)
		code := 500
		want := quickapi.ErrorResponse{Code: code, Error: "internal server error"}

		req := httptest.NewRequest("GET", "/todos/", nil)
		got := quickapitest.DoRequest[quickapi.ErrorResponse](t, req, code, handler)

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("%s %s, response mismatch (-want +got):\n%s", req.Method, req.URL.Path, diff)
		}
	})
}

func TestPostTodo(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		code := 201
		handler := http.HandlerFunc(webapi.Mount(chi.NewRouter(), &MockTodoRepository{}).ServeHTTP)
		payload := strings.NewReader(`{"title":"test-title","content":"test-content"}`)

		req := httptest.NewRequest("POST", "/todos/", payload)
		quickapitest.DoRequest[any](t, req, code, handler, func(t *testing.T, res *http.Response) {
			if res.Header.Get("Location") != req.Host+req.URL.Path+"2" {
				t.Errorf("Location is %v", res.Header.Get("Location"))
			}
		})
	})

	// TODO: bad request

	t.Run("error", func(t *testing.T) {
		handler := http.HandlerFunc(webapi.Mount(chi.NewRouter(), &MockTodoRepositoryError{}).ServeHTTP)
		code := 500
		want := quickapi.ErrorResponse{Code: code, Error: "internal server error"}

		payload := strings.NewReader(`{"title":"test-title","content":"test-content"}`)
		req := httptest.NewRequest("POST", "/todos/", payload)
		got := quickapitest.DoRequest[quickapi.ErrorResponse](t, req, code, handler)

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("%s %s, response mismatch (-want +got):\n%s", req.Method, req.URL.Path, diff)
		}
	})
}

func TestPutTodo(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		handler := http.HandlerFunc(webapi.Mount(chi.NewRouter(), &MockTodoRepository{}).ServeHTTP)
		code := 204

		payload := strings.NewReader(`{"title":"test-title","content":"test-content"}`)
		req := httptest.NewRequest("PUT", "/todos/2", payload)
		quickapitest.DoRequest[any](t, req, code, handler)
	})

	t.Run("invalid-path", func(t *testing.T) {
		handler := http.HandlerFunc(webapi.Mount(chi.NewRouter(), &MockTodoRepository{}).ServeHTTP)
		code := http.StatusMethodNotAllowed

		payload := strings.NewReader(`{"title":"test-title","content":"test-content"}`)
		req := httptest.NewRequest("PUT", "/todos/", payload)
		quickapitest.DoRequest[any](t, req, code, handler)
	})

	t.Run("error", func(t *testing.T) {
		code := 500
		handler := http.HandlerFunc(webapi.Mount(chi.NewRouter(), &MockTodoRepositoryError{}).ServeHTTP)

		payload := strings.NewReader(`{"title":"test-title","content":"test-content"}`)
		req := httptest.NewRequest("PUT", "/todos/2", payload)
		quickapitest.DoRequest[any](t, req, code, handler)
	})
}

func TestDeleteTodo(t *testing.T) {
	handler := http.HandlerFunc(webapi.Mount(chi.NewRouter(), &MockTodoRepository{}).ServeHTTP)
	t.Run("ok", func(t *testing.T) {
		code := 204

		req := httptest.NewRequest("DELETE", "/todos/2", nil)
		quickapitest.DoRequest[any](t, req, code, handler)
	})

	t.Run("invalid-path", func(t *testing.T) {
		code := http.StatusMethodNotAllowed

		req := httptest.NewRequest("DELETE", "/todos/", nil)
		quickapitest.DoRequest[any](t, req, code, handler)
	})

	t.Run("error", func(t *testing.T) {
		handler := http.HandlerFunc(webapi.Mount(chi.NewRouter(), &MockTodoRepositoryError{}).ServeHTTP)
		code := 500

		req := httptest.NewRequest("DELETE", "/todos/2", nil)
		quickapitest.DoRequest[any](t, req, code, handler)
	})
}

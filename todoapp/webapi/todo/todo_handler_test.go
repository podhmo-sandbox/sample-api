package todo_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	webapi "github.com/podhmo-sandbox/sample-api/todoapp/webapi/todo"
	"github.com/podhmo/or"
)

// TODO: performance up

func TestGetTodos(t *testing.T) {
	router := webapi.Mount(chi.NewRouter(), &MockTodoRepository{})
	ts := httptest.NewServer(router)
	defer ts.Close()

	t.Run("not-found", func(t *testing.T) {
		req := or.Fatal(http.NewRequest("GET", ts.URL+"/todos/", nil))(t)
		res := or.Fatal(http.DefaultClient.Do(req))(t)

		if res.StatusCode != 200 {
			t.Errorf("Response code is %v", res.StatusCode)
		}
		if res.Header.Get("Content-Type") != "application/json" {
			t.Errorf("Content-Type is %v", res.Header.Get("Content-Type"))
		}

		var todosResponse webapi.TodosResponse
		if err := json.NewDecoder(res.Body).Decode(&todosResponse); err != nil {
			t.Errorf("unexpected errpr (decode) %+v", err)
		}
		if len(todosResponse.Todos) != 0 {
			t.Errorf("Response is %v", todosResponse.Todos)
		}
	})

	t.Run("ok", func(t *testing.T) {
		router := webapi.Mount(chi.NewRouter(), &MockTodoRepositoryGetTodosExist{})
		ts := httptest.NewServer(router)
		defer ts.Close()

		req := or.Fatal(http.NewRequest("GET", ts.URL+"/todos/", nil))(t)
		res := or.Fatal(http.DefaultClient.Do(req))(t)

		if res.StatusCode != 200 {
			t.Errorf("Response code is %v", res.StatusCode)
		}
		if res.Header.Get("Content-Type") != "application/json" {
			t.Errorf("Content-Type is %v", res.Header.Get("Content-Type"))
		}

		var todosResponse webapi.TodosResponse
		if err := json.NewDecoder(res.Body).Decode(&todosResponse); err != nil {
			t.Errorf("unexpected errpr (decode) %+v", err)
		}
		if len(todosResponse.Todos) != 2 {
			t.Errorf("Response is %v", todosResponse.Todos)
		}
	})

	t.Run("error", func(t *testing.T) {
		router := webapi.Mount(chi.NewRouter(), &MockTodoRepositoryError{})
		ts := httptest.NewServer(router)
		defer ts.Close()

		req := or.Fatal(http.NewRequest("GET", ts.URL+"/todos/", nil))(t)
		res := or.Fatal(http.DefaultClient.Do(req))(t)

		if res.StatusCode != 500 {
			t.Errorf("Response cod is %v", res.StatusCode)
		}
		if res.Header.Get("Content-Type") != "" {
			t.Errorf("Content-Type is %v", res.Header.Get("Content-Type"))
		}

		// if res.Body.Len() != 0 {
		// 	t.Errorf("body is %v", res.Body.Len())
		// }
	})
}

func TestPostTodo(t *testing.T) {
	router := webapi.Mount(chi.NewRouter(), &MockTodoRepository{})
	ts := httptest.NewServer(router)
	defer ts.Close()

	t.Run("ok", func(t *testing.T) {
		payload := strings.NewReader(`{"title":"test-title","content":"test-content"}`)
		req := or.Fatal(http.NewRequest("POST", ts.URL+"/todos/", payload))(t)
		res := or.Fatal(http.DefaultClient.Do(req))(t)

		if res.StatusCode != 201 {
			t.Errorf("Response code is %v", res.StatusCode)
		}
		if res.Header.Get("Location") != req.Host+req.URL.Path+"2" {
			t.Errorf("Location is %v", res.Header.Get("Location"))
		}
	})

	t.Run("error", func(t *testing.T) {
		router := webapi.Mount(chi.NewRouter(), &MockTodoRepositoryError{})
		ts := httptest.NewServer(router)
		defer ts.Close()

		payload := strings.NewReader(`{"title":"test-title","content":"test-content"}`)
		req := or.Fatal(http.NewRequest("POST", ts.URL+"/todos/", payload))(t)
		res := or.Fatal(http.DefaultClient.Do(req))(t)

		if res.StatusCode != 500 {
			t.Errorf("Response code is %v", res.StatusCode)
		}
		if res.Header.Get("Location") != "" {
			t.Errorf("Location is %v", res.Header.Get("Location"))
		}
	})
}

func TestPutTodo(t *testing.T) {
	router := webapi.Mount(chi.NewRouter(), &MockTodoRepository{})
	ts := httptest.NewServer(router)
	defer ts.Close()

	t.Run("ok", func(t *testing.T) {
		req := or.Fatal(http.NewRequest("PUT", ts.URL+"/todos/2", nil))(t)
		res := or.Fatal(http.DefaultClient.Do(req))(t)

		if res.StatusCode != 204 {
			t.Errorf("Response cod is %v", res.StatusCode)
		}
	})

	t.Run("invalid-path", func(t *testing.T) {
		req := or.Fatal(http.NewRequest("PUT", ts.URL+"/todos/", nil))(t)
		res := or.Fatal(http.DefaultClient.Do(req))(t)

		if res.StatusCode != 405 {
			t.Errorf("Response cod is %v", res.StatusCode)
		}
	})

	t.Run("error", func(t *testing.T) {
		router := webapi.Mount(chi.NewRouter(), &MockTodoRepositoryError{})
		ts := httptest.NewServer(router)
		defer ts.Close()

		payload := strings.NewReader(`{"title":"test-title","contents":"test-content"}`)
		req := or.Fatal(http.NewRequest("PUT", ts.URL+"/todos/2", payload))(t)
		res := or.Fatal(http.DefaultClient.Do(req))(t)

		if res.StatusCode != 500 {
			t.Errorf("Response cod is %v", res.StatusCode)
		}
	})

}

func TestDeleteTodo(t *testing.T) {
	router := webapi.Mount(chi.NewRouter(), &MockTodoRepository{})
	ts := httptest.NewServer(router)
	defer ts.Close()

	t.Run("ok", func(t *testing.T) {
		req := or.Fatal(http.NewRequest("DELETE", ts.URL+"/todos/2", nil))(t)
		res := or.Fatal(http.DefaultClient.Do(req))(t)
		if res.StatusCode != 204 {
			t.Errorf("Response cod is %v", res.StatusCode)
		}
	})

	t.Run("invalid-path", func(t *testing.T) {
		req := or.Fatal(http.NewRequest("DELETE", ts.URL+"/todos/", nil))(t)
		res := or.Fatal(http.DefaultClient.Do(req))(t)
		if res.StatusCode != 405 {
			t.Errorf("Response cod is %v", res.StatusCode)
		}
	})

	t.Run("error", func(t *testing.T) {
		router := webapi.Mount(chi.NewRouter(), &MockTodoRepositoryError{})
		ts := httptest.NewServer(router)
		defer ts.Close()

		req := or.Fatal(http.NewRequest("DELETE", ts.URL+"/todos/2", nil))(t)
		res := or.Fatal(http.DefaultClient.Do(req))(t)
		if res.StatusCode != 500 {
			t.Errorf("Response cod is %v", res.StatusCode)
		}
	})
}

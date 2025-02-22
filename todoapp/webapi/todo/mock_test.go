package todo_test

import (
	"errors"

	"github.com/podhmo-sandbox/sample-api/todoapp/entity"
)

type MockTodoRepository struct {
}

func (mtr *MockTodoRepository) GetTodos() (todos []entity.Todo, err error) {
	todos = []entity.Todo{}
	return
}

func (mtr *MockTodoRepository) InsertTodo(todo entity.Todo) (id int, err error) {
	id = 2
	return
}

func (mtr *MockTodoRepository) UpdateTodo(todo entity.Todo) (err error) {
	return
}

func (mtr *MockTodoRepository) DeleteTodo(id int) (err error) {
	return
}

type MockTodoRepositoryGetTodosExist struct {
}

func (mtrgex *MockTodoRepositoryGetTodosExist) GetTodos() (todos []entity.Todo, err error) {
	todos = []entity.Todo{}
	todos = append(todos, entity.Todo{ID: 1, Title: "title1", Content: "contents1"})
	todos = append(todos, entity.Todo{ID: 2, Title: "title2", Content: "contents2"})
	return
}

func (mtrgex *MockTodoRepositoryGetTodosExist) InsertTodo(todo entity.Todo) (id int, err error) {
	return
}

func (mtrgex *MockTodoRepositoryGetTodosExist) UpdateTodo(todo entity.Todo) (err error) {
	return
}

func (mtrgex *MockTodoRepositoryGetTodosExist) DeleteTodo(id int) (err error) {
	return
}

type MockTodoRepositoryError struct {
}

func (mtrgtn *MockTodoRepositoryError) GetTodos() (todos []entity.Todo, err error) {
	err = errors.New("unexpected error occurred")
	return
}

func (mtrgie *MockTodoRepositoryError) InsertTodo(todo entity.Todo) (id int, err error) {
	err = errors.New("unexpected error occurred")
	return
}

func (mtrgue *MockTodoRepositoryError) UpdateTodo(todo entity.Todo) (err error) {
	err = errors.New("unexpected error occurred")
	return
}

func (mtrgde *MockTodoRepositoryError) DeleteTodo(id int) (err error) {
	err = errors.New("unexpected error occurred")
	return
}

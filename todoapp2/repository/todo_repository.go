package repository

import (
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/podhmo-sandbox/sample-api/todoapp2/entity"
)

type TodoRepository struct {
	DB *sqlx.DB
}

func NewTodoRepository(db *sqlx.DB) *TodoRepository {
	return &TodoRepository{DB: db}
}

func (tr *TodoRepository) GetTodos() ([]entity.Todo, error) {
	var todos []entity.Todo
	stmt := "SELECT id, title, content FROM todo ORDER BY id DESC"
	if err := tr.DB.Select(&todos, stmt); err != nil {
		log.Print(err)
		return nil, err
	}
	return todos, nil
}

func (tr *TodoRepository) InsertTodo(todo entity.Todo) (int, error) {
	var id int
	stmt := `INSERT INTO todo (title, content) VALUES (?, ?) RETURNING id`
	if err := tr.DB.Get(&id, stmt, todo.Title, todo.Content); err != nil {
		log.Print(err)
		return id, err
	}
	return id, nil
}

func (tr *TodoRepository) UpdateTodo(todo entity.Todo) error {
	stmt := "UPDATE todo SET title = ?, content = ? WHERE id = ?"
	if _, err := tr.DB.Exec(stmt, todo.Title, todo.Content, todo.ID); err != nil {
		return err
	}
	return nil
}

func (tr *TodoRepository) DeleteTodo(id int) error {
	stmt := "DELETE FROM todo WHERE id = ?"
	if _, err := tr.DB.Exec(stmt, id); err != nil {
		return err
	}
	return nil
}

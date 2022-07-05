package repository

import (
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/podhmo-sandbox/sample-api/model/entity"
)

type TodoRepository interface {
	GetTodos() (todos []entity.Todo, err error)
	InsertTodo(todo entity.Todo) (id int, err error)
	UpdateTodo(todo entity.Todo) (err error)
	DeleteTodo(id int) (err error)
}

// this is temporary implementation
type todoRepository struct {
	DB *sqlx.DB
}

func NewTodoRepository() *todoRepository {
	return &todoRepository{DB: Db}
}

func (tr *todoRepository) GetTodos() (todos []entity.Todo, err error) {
	todos = []entity.Todo{}
	rows, err := tr.DB.
		Query("SELECT id, title, content FROM todo ORDER BY id DESC")
	if err != nil {
		log.Print(err)
		return
	}

	for rows.Next() {
		todo := entity.Todo{}
		err = rows.Scan(&todo.Id, &todo.Title, &todo.Content)
		if err != nil {
			log.Print(err)
			return
		}
		todos = append(todos, todo)
	}

	return
}

func (tr *todoRepository) InsertTodo(todo entity.Todo) (int, error) {
	var id int
	stmt := `INSERT INTO todo (title, content) VALUES (?, ?) RETURNING id`
	if err := tr.DB.Get(&id, stmt, todo.Title, todo.Content); err != nil {
		log.Print(err)
		return id, err
	}
	return id, nil
}

func (tr *todoRepository) UpdateTodo(todo entity.Todo) (err error) {
	_, err = tr.DB.Exec("UPDATE todo SET title = ?, content = ? WHERE id = ?", todo.Title, todo.Content, todo.Id)
	return
}

func (tr *todoRepository) DeleteTodo(id int) (err error) {
	_, err = tr.DB.Exec("DELETE FROM todo WHERE id = ?", id)
	return
}

package entity

type TodoEntity struct {
	Id      int    `db:"id"`
	Title   string `db:"title"`
	Content string `db:"content"`
}

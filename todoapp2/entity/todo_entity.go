package entity

type Todo struct {
	ID      int    `db:"id"`
	Title   string `db:"title"`
	Content string `db:"content"`
}

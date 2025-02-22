package repository

import (
	"context"
	"sort"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/podhmo-sandbox/sample-api/pkg/dblib"
	"github.com/podhmo-sandbox/sample-api/todoapp/entity"
	"github.com/podhmo-sandbox/sample-api/todoapp/repository/repositorytest"
	_ "modernc.org/sqlite"
)

func TestGetTodos(t *testing.T) {
	ctx := context.Background()

	todos := []entity.Todo{
		{ID: 10, Title: "go to bed", Content: "should sleep"},
		{ID: 11, Title: "go to toilet", Content: "should"},
	}
	db, teardown := dblib.NewDB(ctx, t, repositorytest.WithTodo(todos))
	defer teardown()

	dblib.AssertRowsCount(t, db, "todo", 2 /* want*/) // todo: checking by defer

	repo := &TodoRepository{DB: db}
	got, err := repo.GetTodos()
	if err != nil {
		t.Errorf("unexpected error: %+v", err)
	}

	// order by id desc
	want := todos
	sort.Slice(want, func(i, j int) bool { return want[i].ID > want[j].ID })

	type ref struct{ XS []entity.Todo }
	if diff := cmp.Diff(ref{want}, ref{got}); diff != "" {
		t.Errorf("GetContext() mismatch (-want +got):\n%s", diff)
	}
}

func TestInsertTodo(t *testing.T) {
	ctx := context.Background()
	db, teardown := dblib.NewDB(ctx, t, repositorytest.WithTodo(nil))
	defer teardown()

	assertAfterAction := dblib.AssertRowsCountWith(t, db, "todo", 0 /* want*/)
	defer assertAfterAction(1 /* want */)

	want := entity.Todo{
		Title:   "go to bed",
		Content: "should sleep",
	}

	repo := &TodoRepository{DB: db}
	id, err := repo.InsertTodo(want)
	if err != nil {
		t.Errorf("unexpected error: %+v", err)
	}

	var got entity.Todo
	if err := db.GetContext(ctx, &got, "SELECT id,title,content FROM todo WHERE id=?", id); err != nil {
		t.Errorf("unexpected error (db check): %+v", err)
	}
	if diff := cmp.Diff(want, got, cmpopts.IgnoreFields(entity.Todo{}, "ID")); diff != "" {
		t.Errorf("GetContext() mismatch (-want +got):\n%s", diff)
	}
}

func TestUpdateTodo(t *testing.T) {
	ctx := context.Background()

	id := 10
	todos := []entity.Todo{
		{ID: id, Title: "go to bed", Content: "should sleep"},
	}
	db, teardown := dblib.NewDB(ctx, t, repositorytest.WithTodo(todos))
	defer teardown()
	dblib.AssertRowsCount(t, db, "todo", 1 /* want*/)

	want := entity.Todo{ID: id, Title: "*", Content: "**"}
	repo := &TodoRepository{DB: db}
	if err := repo.UpdateTodo(want); err != nil {
		t.Errorf("unexpected error: %+v", err)
	}

	gots, err := repo.GetTodos()
	if err != nil {
		t.Errorf("unexpected error: %+v", err)
	}
	if diff := cmp.Diff(want, gots[0]); diff != "" {
		t.Errorf("GetContext() mismatch (-want +got):\n%s", diff)
	}
}

func TestDeleteTodo(t *testing.T) {
	ctx := context.Background()
	id := 10
	todos := []entity.Todo{
		{ID: id, Title: "go to bed", Content: "should sleep"},
	}
	db, teardown := dblib.NewDB(ctx, t, repositorytest.WithTodo(todos))
	defer teardown()

	assertAfterAction := dblib.AssertRowsCountWith(t, db, "todo", 1 /* want*/)
	defer assertAfterAction(0 /* want */)

	repo := &TodoRepository{DB: db}
	if err := repo.DeleteTodo(todos[0].ID); err != nil {
		t.Errorf("unexpected error: %+v", err)
	}
}

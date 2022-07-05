package repository

import (
	"context"
	"sort"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/podhmo-sandbox/sample-api/model/entity"
	rt "github.com/podhmo-sandbox/sample-api/model/repository/repositorytest"
)

func TestGetTodos(t *testing.T) {
	ctx := context.Background()

	todos := []entity.Todo{
		{Id: 10, Title: "go to bed", Content: "should sleep"},
		{Id: 11, Title: "go to toilet", Content: "should"},
	}
	db, teardown := rt.NewDB(ctx, t, rt.WithTodo(todos))
	defer teardown()

	rt.AssertRowsCount(t, db, "todo", 2 /* want*/) // todo: checking by defer

	repo := &todoRepository{DB: db}
	got, err := repo.GetTodos()
	if err != nil {
		t.Errorf("unexpected error: %+v", err)
	}

	// order by id desc
	want := todos
	sort.Slice(want, func(i, j int) bool { return want[i].Id > want[j].Id })

	type ref struct{ XS []entity.Todo }
	if diff := cmp.Diff(ref{want}, ref{got}); diff != "" {
		t.Errorf("GetContext() mismatch (-want +got):\n%s", diff)
	}
}

func TestInsertTodo(t *testing.T) {
	ctx := context.Background()
	db, teardown := rt.NewDB(ctx, t, rt.WithTodo(nil))
	defer teardown()

	rt.AssertRowsCount(t, db, "todo", 0 /* want*/) // todo: checking by defer
	want := entity.Todo{
		Title:   "go to bed",
		Content: "should sleep",
	}

	repo := &todoRepository{DB: db}
	id, err := repo.InsertTodo(want)
	if err != nil {
		t.Errorf("unexpected error: %+v", err)
	}

	var got entity.Todo
	if err := db.GetContext(ctx, &got, "SELECT id,title,content FROM todo WHERE id=?", id); err != nil {
		t.Errorf("unexpected error (db check): %+v", err)
	}
	if diff := cmp.Diff(want, got, cmpopts.IgnoreFields(entity.Todo{}, "Id")); diff != "" {
		t.Errorf("GetContext() mismatch (-want +got):\n%s", diff)
	}
	rt.AssertRowsCount(t, db, "todo", 1 /* want*/)
}

func TestUpdateTodo(t *testing.T) {
	ctx := context.Background()

	id := 10
	todos := []entity.Todo{
		{Id: id, Title: "go to bed", Content: "should sleep"},
	}
	db, teardown := rt.NewDB(ctx, t, rt.WithTodo(todos))
	defer teardown()
	rt.AssertRowsCount(t, db, "todo", 1 /* want*/)

	want := entity.Todo{Id: id, Title: "*", Content: "**"}
	repo := &todoRepository{DB: db}
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

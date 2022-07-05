package repository

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/podhmo-sandbox/sample-api/model/entity"
	rt "github.com/podhmo-sandbox/sample-api/model/repository/repositorytest"
)

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

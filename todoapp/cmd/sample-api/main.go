package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jmoiron/sqlx"
	"github.com/podhmo-sandbox/sample-api/pkg/dblib"
	"github.com/podhmo-sandbox/sample-api/todoapp/repository"
	"github.com/podhmo-sandbox/sample-api/todoapp/webapi/todo"
	"github.com/podhmo/flagstruct"
	_ "modernc.org/sqlite"
)

type Config struct {
	DB   dblib.Config `flag:"db"`
	Addr string       `flag:"addr"`
}

func main() {
	config := &Config{DB: dblib.DefaultConfig(), Addr: ":8080"} // default values
	flagstruct.Parse(config)
	if err := run(*config); err != nil {
		log.Fatalf("!! %+v", err)
	}
}

func run(config Config) error {
	// json.NewEncoder(os.Stdout).Encode(config)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	server := http.Server{
		Addr: config.Addr,
	}
	go func() {
		<-ctx.Done()
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		server.Shutdown(ctx)
	}()

	r := chi.NewRouter()

	// A good base middleware stack
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	db, err := config.DB.New(ctx)
	if err != nil {
		return err
	}
	mount(r, db)
	return server.ListenAndServe()
}

func mount(r chi.Router, db *sqlx.DB) {
	{
		repo := repository.NewTodoRepository(db)
		todo.Mount(r, repo)
	}
}

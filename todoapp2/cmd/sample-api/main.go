package main

import (
	"context"
	"log"
	"time"

	"github.com/podhmo-sandbox/sample-api/pkg/dblib"
	"github.com/podhmo-sandbox/sample-api/todoapp2/repository"
	"github.com/podhmo-sandbox/sample-api/todoapp2/webapi/todo"
	"github.com/podhmo/flagstruct"
	"github.com/podhmo/quickapi"
	"github.com/podhmo/quickapi/experimental/define"
	rohandler "github.com/podhmo/reflect-openapi/handler"
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
	ctx := context.Background()

	db, err := config.DB.New(ctx)
	if err != nil {
		return err
	}

	r := quickapi.DefaultRouter()
	bc, err := define.NewBuildContext(define.Doc(), r)
	if err != nil {
		return err
	}

	{
		repo := repository.NewTodoRepository(db)
		todo.Mount(repo)(bc)
	}

	handler, err := bc.BuildHandler(ctx)
	if err != nil {
		return err
	}
	bc.Router().Mount("/openapi", rohandler.NewHandler(bc.Doc(), "/openapi"))
	return quickapi.NewServer(config.Addr, handler, 5*time.Second).ListenAndServe(ctx)
}

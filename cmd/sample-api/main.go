package main

import (
	"net/http"

	"github.com/podhmo-sandbox/sample-api/controller"
	"github.com/podhmo-sandbox/sample-api/controller/router"
	"github.com/podhmo-sandbox/sample-api/model/repository"
)

var tr = repository.NewTodoRepository()
var tc = controller.NewTodoController(tr)
var ro = router.NewRouter(tc)

func main() {
	server := http.Server{
		Addr: ":8080",
	}
	http.HandleFunc("/todos/", ro.HandleTodosRequest)
	server.ListenAndServe()
}

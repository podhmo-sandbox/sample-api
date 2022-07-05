package main

import (
	"net/http"

	"github.com/podhmo-sandbox/sample-api/controller"
	"github.com/podhmo-sandbox/sample-api/controller/router"
	"github.com/podhmo-sandbox/sample-api/model/repository"
)

type Router interface {
	HandleTodosRequest(w http.ResponseWriter, r *http.Request)
}

func mount(ro Router) {
	http.HandleFunc("/todos/", ro.HandleTodosRequest)
}

func main() {
	server := http.Server{
		Addr: ":8080",
	}

	tr := repository.NewTodoRepository()
	tc := controller.NewTodoController(tr)
	ro := router.NewRouter(tc)
	mount(ro)
	server.ListenAndServe()
}

package main

import (
	"github.com/lexizz/metcoll/cmd/server/handlers"
	"log"
	"net/http"
)

const (
	HOST string = "127.0.0.1"
	PORT string = "8080"
)

type server struct {
	httpServer *http.Server
}

func New() *server {
	return &server{
		httpServer: &http.Server{
			Addr:    HOST + ":" + PORT,
			Handler: getRoutes(),
		},
	}
}

func getRoutes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.Handle("/", http.NotFoundHandler())
	mux.HandleFunc("/update/", handlers.UpdateMetric())

	return mux
}

func main() {
	serv := New()

	log.Fatal(serv.httpServer.ListenAndServe())
}

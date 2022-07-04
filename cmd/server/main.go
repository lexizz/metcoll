package main

import (
	"log"
	"net/http"
)

const (
	HOST string = "127.0.0.1"
	PORT string = "8080"
)

func getRoutes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.Handle("/", http.NotFoundHandler())
	mux.HandleFunc("/update/", POSTHandler)

	return mux
}

func main() {
	server := &http.Server{
		Addr:    HOST + ":" + PORT,
		Handler: getRoutes(),
	}

	log.Fatal(server.ListenAndServe())
}

func POSTHandler(w http.ResponseWriter, r *http.Request) {
	log.Println(r.Method, r.Host, r.URL.Path, r.URL.RawQuery)

	if r.Method != http.MethodPost {
		http.Error(w, "Only POST requests are allowed!", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("content-type", "text/plain")

	w.WriteHeader(http.StatusOK)

	_, writeError := w.Write([]byte("Ok"))
	if writeError != nil {
		log.Fatal(writeError)
	}
}

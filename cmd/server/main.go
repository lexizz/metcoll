package main

import (
	"log"

	"github.com/lexizz/metcoll/internal/server"
)

func main() {
	serv := server.New()

	log.Fatal(serv.Init.ListenAndServe())
}

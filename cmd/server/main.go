package main

import (
	"log"

	"github.com/lexizz/metcoll/internal/server"
)

func main() {
	httpServer := server.New()

	log.Fatal(httpServer.Listen())
}

package main

import (
	"github.com/joostvdg/proglog/internal/server"
	"log"
)

func main() {

	srv := server.NewHTTPServer(":8081")
	log.Fatal(srv.ListenAndServe())
}

package main

import (
	"database/sql"
	"log"
	"net/http"

	"os"

	"github.com/ob-vss-ss18/ppl-auth/api"
	"github.com/ob-vss-ss18/ppl-auth/backend"
)

var db *sql.DB

func main() {
	log.Printf("Initializing...\n")

	backend.ConnectDb()
	backend.MigrateDb()

	http.Handle("/graphql", api.ApiHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Starting listener on port %s\n", port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal(err)
	}
}

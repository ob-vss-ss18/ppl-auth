package main

import (
	"database/sql"
	"log"
	"net/http"

	"os"

	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/source/file"
	_ "github.com/lib/pq"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
)

var db *sql.DB

func main() {
	log.Printf("Initializing...\n")
	log.Printf("Connecting to database '%s'\n", os.Getenv("DATABASE_URL"))
	db, _ = sql.Open("postgres", os.Getenv("DATABASE_URL"))
	//if err != nil {
	//	log.Fatal(err)
	//}
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatal(err)
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations/",
		"postgres", driver)
	if err != nil {
		log.Fatal(err)
	}
	m.Log = migrationLogger{false}
	log.Printf("Running needed migrations\n")
	m.Up()

	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query: QueryType,
		// Mutation: MutationType,
	})
	if err != nil {
		log.Fatal(err)
	}

	handler := handler.New(&handler.Config{
		Schema:   &schema,
		GraphiQL: true,
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	http.Handle("/graphql", handler)
	log.Printf("Starting listener on port %s\n", port)
	err = http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal(err)
	}
}

type migrationLogger struct {
	verbose bool
}

func (m migrationLogger) Printf(format string, v ...interface{}) {
	log.Printf(format, v...)
}

func (m migrationLogger) Verbose() bool {
	return m.verbose
}

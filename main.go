package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"os"

	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/source/file"
	_ "github.com/lib/pq"

	"github.com/graphql-go/graphql"
)

var db *sql.DB

func main() {
	log.Printf("Initializing...\n")
	log.Printf("Connecting to database '%s'\n", os.Getenv("DATABASE_URL"))
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}
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
		Query:    QueryType,
		Mutation: MutationType,
	})
	if err != nil {
		log.Fatal(err)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	http.HandleFunc("/", hello)
	http.Handle("/graphql", handler(schema))
	log.Printf("Starting listener on port %s\n", port)
	err = http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal(err)
	}
}

func hello(res http.ResponseWriter, req *http.Request) {
	fmt.Fprintln(res, "Hello World!")
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

func handler(schema graphql.Schema) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		query, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		result := graphql.Do(graphql.Params{
			Schema:        schema,
			RequestString: string(query),
		})
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(result)
	}
}

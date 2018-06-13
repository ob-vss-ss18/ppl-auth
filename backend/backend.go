package backend

import (
	"database/sql"
	"log"
	"os"

	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/source/file"
	_ "github.com/lib/pq"
)

var db *sql.DB

func ConnectDb() {
	var err error
	log.Printf("Connecting to database '%s'\n", os.Getenv("DATABASE_URL"))
	db, err = sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}
}

func MigrateDb() {
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

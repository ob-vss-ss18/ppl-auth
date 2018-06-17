package api

import (
	"log"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
)

var ApiHandler *handler.Handler

func init() {
	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query:    QueryType,
		Mutation: MutationType,
	})
	if err != nil {
		log.Fatal(err)
	}

	ApiHandler = handler.New(&handler.Config{
		Schema:   &schema,
		GraphiQL: true,
	})
}

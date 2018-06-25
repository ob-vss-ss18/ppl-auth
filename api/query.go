package api

import (
	"github.com/graphql-go/graphql"
	"github.com/ob-vss-ss18/ppl-auth/backend"
)

var QueryType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Query",
	Fields: graphql.Fields{
		"user": &graphql.Field{
			Type: UserType,
			Args: graphql.FieldConfigArgument{
				"email": &graphql.ArgumentConfig{
					Description: "E-Mail address",
					Type:        graphql.NewNonNull(graphql.String),
				},
				"token": &graphql.ArgumentConfig{
					Description: "Token",
					Type:        graphql.NewNonNull(graphql.String),
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				email := p.Args["email"].(string)
				token := p.Args["token"].(string)
				return backend.ValidateToken(email, token)
			},
		},
	},
})

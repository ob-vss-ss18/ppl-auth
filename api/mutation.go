package api

import (
	"github.com/graphql-go/graphql"
	"github.com/ob-vss-ss18/ppl-auth/backend"
)

var MutationType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Mutations",
	Fields: graphql.Fields{
		"loginPwd": &graphql.Field{
			Type: UserType,
			Args: graphql.FieldConfigArgument{
				"email": &graphql.ArgumentConfig{
					Description: "E-Mail address",
					Type:        graphql.NewNonNull(graphql.String),
				},
				"password": &graphql.ArgumentConfig{
					Description: "Password",
					Type:        graphql.NewNonNull(graphql.String),
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				email := p.Args["email"].(string)
				passwd := p.Args["password"].(string)
				return backend.LoginPwd(email, passwd)
			},
		},
		"requestToken": &graphql.Field{
			Type: graphql.Boolean,
			Args: graphql.FieldConfigArgument{
				"email": &graphql.ArgumentConfig{
					Description: "E-Mail address",
					Type:        graphql.NewNonNull(graphql.String),
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				email := p.Args["email"].(string)
				return backend.RequestToken(email)
			},
		},
	},
})

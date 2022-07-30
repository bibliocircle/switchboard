package gql

import "github.com/graphql-go/graphql"

type HTTPHeader struct {
	Name  string `json:"name"`
	Value string ` json:"value"`
}

var HTTPHeaderGqlType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Header",
	Fields: graphql.Fields{
		"name": &graphql.Field{
			Type: graphql.String,
		},
		"value": &graphql.Field{
			Type: graphql.String,
		},
	},
})

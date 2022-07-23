package gql

import "github.com/graphql-go/graphql"

type HeaderConfig struct {
	Name  string `json:"name"`
	Value string ` json:"value"`
}

var HeaderConfigGqlType = graphql.NewObject(graphql.ObjectConfig{
	Name: "KeyValuePair",
	Fields: graphql.Fields{
		"name": &graphql.Field{
			Type: graphql.String,
		},
		"value": &graphql.Field{
			Type: graphql.String,
		},
	},
})

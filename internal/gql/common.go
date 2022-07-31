package gql

import "github.com/graphql-go/graphql"

var HTTPHeaderGqlInputType = graphql.NewInputObject(graphql.InputObjectConfig{
	Name: "HeaderInput",
	Fields: graphql.InputObjectConfigFieldMap{
		"name": &graphql.InputObjectFieldConfig{
			Type: graphql.NewNonNull(graphql.String),
		},
		"value": &graphql.InputObjectFieldConfig{
			Type: graphql.NewNonNull(graphql.String),
		},
	},
})

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

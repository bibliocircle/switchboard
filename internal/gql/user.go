package gql

import (
	"switchboard/internal/db"
	"switchboard/internal/models"

	"github.com/graphql-go/graphql"
)

var UserGqlType = graphql.NewObject(graphql.ObjectConfig{
	Name: "User",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type: graphql.String,
		},
		"firstName": &graphql.Field{
			Type: graphql.String,
		},
		"lastName": &graphql.Field{
			Type: graphql.String,
		},
		"email": &graphql.Field{
			Type: graphql.String,
		},
		"createdAt": &graphql.Field{
			Type: graphql.String,
		},
		"updatedAt": &graphql.Field{
			Type: graphql.String,
		},
	},
})

var UserResolver = func(p graphql.ResolveParams) (interface{}, error) {
	userId, ok := p.Args["id"].(string)
	if ok {
		user, err := db.GetUserByID(userId)
		if err != nil {
			return nil, err
		}
		return user, nil
	}
	return nil, nil
}

var UsersResolver = func(p graphql.ResolveParams) (interface{}, error) {
	users, err := db.GetUsers()
	if err != nil {
		return make([]models.User, 0), err
	}
	return users, nil
}

package user

import (
	"os/user"
	"switchboard/internal/common"
	"switchboard/internal/gql"

	"github.com/graphql-go/graphql"
	"github.com/sirupsen/logrus"
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

func GetUserResolver(p graphql.ResolveParams) (interface{}, error) {
	userId, ok := p.Args["id"].(string)
	if ok {
		user, err := GetUserByID(userId)
		if err != nil {
			logrus.Errorln(err)
			return nil, gql.NewGqlError(common.ErrorGeneric, "could not retrieve user")
		}
		return user, nil
	}
	return nil, nil
}

func GetUsersResolver(p graphql.ResolveParams) (interface{}, error) {
	users, err := GetUsers()
	if err != nil {
		logrus.Errorln(err)
		return make([]user.User, 0), gql.NewGqlError(common.ErrorGeneric, "could not retrieve users")
	}
	return users, nil
}

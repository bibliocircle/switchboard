package scenario

import (
	"switchboard/internal/common"
	"switchboard/internal/gql"
	"switchboard/internal/upstream"
	"switchboard/internal/user"

	"github.com/graphql-go/graphql"
	"github.com/mitchellh/mapstructure"
	"github.com/sirupsen/logrus"
)

var HTTPResponseScenarioConfigGqlInputType = graphql.NewInputObject(graphql.InputObjectConfig{
	Name: "HTTPResponseScenarioConfigInput",
	Fields: graphql.InputObjectConfigFieldMap{
		"statusCode": &graphql.InputObjectFieldConfig{
			Type: graphql.NewNonNull(graphql.Int),
		},
		"responseBodyTemplate": &graphql.InputObjectFieldConfig{
			Type: graphql.String,
		},
		"responseHeadersTemplate": &graphql.InputObjectFieldConfig{
			Type: graphql.String,
		},
	},
})

var HTTPResponseScenarioConfigGqlType = graphql.NewObject(graphql.ObjectConfig{
	Name: "HTTPResponseScenarioConfig",
	Fields: graphql.Fields{
		"statusCode": &graphql.Field{
			Type: graphql.Int,
		},
		"responseBodyTemplate": &graphql.Field{
			Type: graphql.String,
		},
		"responseHeadersTemplate": &graphql.Field{
			Type: graphql.String,
		},
	},
})

var ProxyScenarioConfigGqlInputType = graphql.NewInputObject(graphql.InputObjectConfig{
	Name: "ProxyScenarioConfigInput",
	Fields: graphql.InputObjectConfigFieldMap{
		"upstreamId": &graphql.InputObjectFieldConfig{
			Type: graphql.NewNonNull(graphql.String),
		},
		"injectHeaders": &graphql.InputObjectFieldConfig{
			Type: graphql.NewList(gql.HTTPHeaderGqlInputType),
		},
	},
})

var ProxyScenarioConfigGqlType = graphql.NewObject(graphql.ObjectConfig{
	Name: "ProxyScenarioConfig",
	Fields: graphql.Fields{
		"upstream": &graphql.Field{
			Type: upstream.UpstreamGqlType,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				upstreamID := p.Source.(*ProxyScenarioConfig).UpstreamID
				if upstreamID == "" {
					return nil, nil
				}
				upstream, err := upstream.GetUpstreamByID(upstreamID)
				if err != nil {
					logrus.Errorln(err)
					return nil, gql.NewGqlError(common.ErrorGeneric, "could not retrieve upstream")
				}
				return *upstream, nil
			},
		},
		"injectHeaders": &graphql.Field{
			Type: graphql.NewList(gql.HTTPHeaderGqlType),
		},
	},
})

var NetworkScenarioConfigGqlInputType = graphql.NewInputObject(graphql.InputObjectConfig{
	Name: "NetworkScenarioConfigInput",
	Fields: graphql.InputObjectConfigFieldMap{
		"type": &graphql.InputObjectFieldConfig{
			Type: graphql.NewNonNull(graphql.String),
		},
	},
})

var NetworkScenarioConfigGqlType = graphql.NewObject(graphql.ObjectConfig{
	Name: "NetworkScenarioConfig",
	Fields: graphql.Fields{
		"type": &graphql.Field{
			Type: graphql.String,
		},
	},
})

var ScenarioGqlInputType = graphql.NewInputObject(graphql.InputObjectConfig{
	Name: "ScenarioInput",
	Fields: graphql.InputObjectConfigFieldMap{
		"endpointId": &graphql.InputObjectFieldConfig{
			Type: graphql.NewNonNull(graphql.String),
		},
		"type": &graphql.InputObjectFieldConfig{
			Type: graphql.NewNonNull(graphql.String),
		},
		"httpResponseScenarioConfig": &graphql.InputObjectFieldConfig{
			Type: HTTPResponseScenarioConfigGqlInputType,
		},
		"proxyScenarioConfig": &graphql.InputObjectFieldConfig{
			Type: ProxyScenarioConfigGqlInputType,
		},
		"networkScenarioConfig": &graphql.InputObjectFieldConfig{
			Type: NetworkScenarioConfigGqlInputType,
		},
	},
})

var ScenarioGqlType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Scenario",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type: graphql.String,
		},
		"type": &graphql.Field{
			Type: graphql.String,
		},
		"httpResponseScenarioConfig": &graphql.Field{
			Type: HTTPResponseScenarioConfigGqlType,
		},
		"proxyScenarioConfig": &graphql.Field{
			Type: ProxyScenarioConfigGqlType,
		},
		"networkScenarioConfig": &graphql.Field{
			Type: NetworkScenarioConfigGqlType,
		},
		"createdBy": &graphql.Field{
			Type: user.UserGqlType,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				userId := p.Source.(Scenario).CreatedBy
				users, err := user.GetUserByID(userId)
				if err != nil {
					logrus.Errorln(err)
					return make([]user.User, 0), gql.NewGqlError(common.ErrorGeneric, "could not resolve createdBy field")
				}
				return users, nil
			},
		},
		"createdAt": &graphql.Field{
			Type: graphql.String,
		},
		"updatedAt": &graphql.Field{
			Type: graphql.String,
		},
	},
})

func CreateScenarioResolver(p graphql.ResolveParams) (interface{}, error) {
	var input CreateScenarioRequestBody
	mapstructure.Decode(p.Args["scenario"], &input)
	currentUser := p.Context.Value(common.REQ_USER_KEY).(*user.User)

	createdScenario, createErr := CreateScenario(currentUser.ID, &input)
	if createErr == nil {
		return createdScenario, nil
	}

	if createErr.ErrorCode == common.ErrorDuplicateEntity {
		return nil, gql.NewGqlError(common.ErrorDuplicateEntity, "duplicate scenario")
	}

	return nil, gql.NewGqlError(common.ErrorGeneric, "generic error")
}

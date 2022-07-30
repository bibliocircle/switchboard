package gql

import "errors"

type GQLError struct {
	error
	extensions map[string]interface{}
}

func (e GQLError) Extensions() map[string]interface{} {
	return e.extensions
}

func NewGqlError(code string, description string) *GQLError {
	return &GQLError{
		error: errors.New(description),
		extensions: map[string]interface{}{
			"code": code,
		},
	}
}

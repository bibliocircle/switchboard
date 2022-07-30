package gql

import "errors"

const GqlDuplicate = "DUPLICATE"
const GqlForbidden = "FORBIDDEN"
const GqlInternalError = "INTERNAL_ERROR"
const GqlNotFound = "NOT_FOUND"
const GqlUnauthorised = "UNAUTHORISED"

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

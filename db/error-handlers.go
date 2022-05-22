package db

import (
	"errors"
	"switchboard/util"

	"go.mongodb.org/mongo-driver/mongo"
)

func WrapDBErrorIfNecessary(err error) *util.DetailedError {
	var e mongo.WriteException
	if errors.As(err, &e) {
		for _, writeErr := range e.WriteErrors {
			if writeErr.Code == 11000 {
				return util.NewDetailedError(util.ErrorDuplicateEntity, "A user with the same email already exists")
			}
		}
	}
	return util.WrapAsDetailedError(err)
}

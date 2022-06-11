package db

import (
	"errors"
	"switchboard/internal/common"

	"go.mongodb.org/mongo-driver/mongo"
)

func WrapDBErrorIfNecessary(err error) *common.DetailedError {
	var e mongo.WriteException
	if errors.As(err, &e) {
		for _, writeErr := range e.WriteErrors {
			if writeErr.Code == 11000 {
				return common.NewDetailedError(common.ErrorDuplicateEntity, "A user with the same email already exists")
			}
		}
	}
	return common.WrapAsDetailedError(err)
}

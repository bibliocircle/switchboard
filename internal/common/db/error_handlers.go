package db

import (
	e "errors"
	"switchboard/internal/common/err_utils"

	"go.mongodb.org/mongo-driver/mongo"
)

func GetDbError(err error) *err_utils.DetailedError {
	var we mongo.WriteException
	if e.As(err, &we) {
		for _, writeErr := range we.WriteErrors {
			if writeErr.Code == 11000 {
				return err_utils.NewDetailedError(err_utils.ErrorDuplicateEntity, "duplicate document exists")
			}
		}
	}
	return err_utils.WrapAsDetailedError(err)
}

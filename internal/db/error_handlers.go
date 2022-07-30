package db

import (
	e "errors"
	"switchboard/internal/common"

	"go.mongodb.org/mongo-driver/mongo"
)

func GetDbError(err error) *common.DetailedError {
	var we mongo.WriteException
	switch {
	case e.As(err, &we):
		for _, writeErr := range we.WriteErrors {
			if writeErr.Code == 11000 {
				return common.NewDetailedError(common.ErrorDuplicateEntity, "duplicate document exists")
			}
		}
	}
	return common.WrapAsDetailedError(err)
}

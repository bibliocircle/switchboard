package util

import "github.com/google/uuid"

func UUIDv4() string {
	id, _ := uuid.NewRandom()
	return id.String()
}

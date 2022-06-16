package randomdata

import (
	sid "github.com/teris-io/shortid"
)

func GetShortId() string {
	id, _ := sid.Generate()
	return id
}

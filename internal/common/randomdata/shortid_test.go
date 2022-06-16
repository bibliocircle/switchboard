package randomdata_test

import (
	"fmt"
	"switchboard/internal/common/randomdata"
	"switchboard/internal/testutils"
	"testing"
)

func TestShortIdGen(t *testing.T) {
	id := randomdata.GetShortId()
	testutils.Equals(t, fmt.Sprintf("%T", id), "string")
}

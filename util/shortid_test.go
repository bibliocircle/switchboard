package util_test

import (
	"fmt"
	"switchboard/testutils"
	"switchboard/util"
	"testing"
)

func TestShortIdGen(t *testing.T) {
	id := util.GetShortId()
	testutils.Equals(t, fmt.Sprintf("%T", id), "string")
}

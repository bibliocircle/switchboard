package common_test

import (
	"fmt"
	"switchboard/internal/common"
	"switchboard/internal/testutils"
	"testing"
)

func TestShortIdGen(t *testing.T) {
	id := common.GetShortId()
	testutils.Equals(t, fmt.Sprintf("%T", id), "string")
}

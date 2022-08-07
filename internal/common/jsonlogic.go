package common

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	jl "github.com/diegoholiveira/jsonlogic/v3"
	"github.com/tidwall/gjson"
)

func lookupJson(expr, input interface{}) interface{} {
	args := expr.([]interface{})
	reqField := args[0].(string)
	lookupExpr := args[1].(string)

	fmt.Println(reqField)
	fmt.Println(lookupExpr)

	ij, err := json.Marshal(input)
	jsonString := string(ij)
	if err != nil {
		return input
	}
	reqFieldValue := gjson.Get(jsonString, reqField).String()
	result := gjson.Get(reqFieldValue, lookupExpr)
	switch result.Type {
	case 0:
		return nil
	case 1:
		return false
	case 2:
		return result.Int()
	case 4:
		return true
	case 3:
		fallthrough
	case 5:
		return result.String()
	}
	return nil
}

func ApplyJsonLogic(jsonlogic, inputData string) bool {
	jl.AddOperator("lookupJson", lookupJson)
	logic := strings.NewReader(jsonlogic)
	data := strings.NewReader(inputData)

	var result bytes.Buffer

	jl.Apply(logic, data, &result)
	return result.String() == "true\n"
}

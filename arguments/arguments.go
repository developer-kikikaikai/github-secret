package arguments

import (
	"flag"
	"fmt"
)

const (
	invalidString = "Please set the value"
)

type argument interface {
	HasArg() bool
}

// string definition
type stringArg struct {
	Val *string
}

func (arg *stringArg) HasArg() bool {
	return (arg.Val != nil && *arg.Val != invalidString)
}

//input param map
var requiredTypeParams = map[string]argument{}

func RequiredString(name string, usage string) *string {
	res := flag.String(name, invalidString, "[required] "+usage)
	requiredTypeParams[name] = &stringArg{Val: res}
	return res
}

func String(name string, value string, usage string) *string {
	return flag.String(name, value, usage)
}

func Parse() bool {
	flag.Parse()
	required := ""
	for name, input := range requiredTypeParams {
		if !input.HasArg() {
			required += name + " "
		}
	}
	if required != "" {
		fmt.Printf("Required fields are not enough, fields: %s\n", required)
		return false
	}
	return true
}

func Usage() {
	flag.Usage()
}

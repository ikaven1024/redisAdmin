package util

import "fmt"

func MustString(o interface{}) string {
	s, ok := o.(string)
	if !ok {
		panic(fmt.Sprintf("%v is not a string.", o))
	}
	return s
}

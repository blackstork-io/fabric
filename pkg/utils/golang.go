package utils

import "reflect"

func IsNil(c any) bool {
	// Checking for nil on interface objects is terrble
	return c == nil || (reflect.ValueOf(c).Kind() == reflect.Ptr && reflect.ValueOf(c).IsNil())
}

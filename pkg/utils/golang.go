package utils

import "reflect"

// Correct version of nil check, works on nil interfaces as well as any other value.
func IsNil(val any) bool {
	// Checking for nil on interface objects is terrble
	// Thanks to: https://stackoverflow.com/a/76595928/4632951
	if val == nil {
		return true
	}
	v := reflect.ValueOf(val)
	k := v.Kind()
	switch k {
	case reflect.Chan, reflect.Func, reflect.Map, reflect.Pointer,
		reflect.UnsafePointer, reflect.Interface, reflect.Slice:
		return v.IsNil()
	}

	return false
}

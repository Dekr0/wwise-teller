package assert

import (
	"fmt"
	"reflect"
)

func Equal[T comparable](expected T, received T, msg string) {
	if expected != received {
		msg = fmt.Sprintf("%s. Expected: %v. Receive: %v", msg, expected, received)
		panic(msg)
	}
}

func True(expr bool, msg string) {
	if !expr {
		panic(msg)
	}
}

func Nil(target any, name string) {
	v := reflect.ValueOf(target)
	if !v.IsValid() {
		return
	}
	switch v.Kind() {
	case reflect.Ptr, reflect.Slice, reflect.Map, reflect.Func, reflect.Interface:
		if v.IsNil() {
			return
		}
		panic(fmt.Sprintf("%s is not nil", name))
	default:
		panic(fmt.Sprintf("%s is not nil", name))
	}
}

func NotNil(target any, name string) {
	v := reflect.ValueOf(target)
	if !v.IsValid() {
		panic(fmt.Sprintf("%s is nil", name))
	}
	switch v.Kind() {
	case reflect.Ptr, reflect.Slice, reflect.Map, reflect.Func, reflect.Interface:
		if v.IsNil() {
			panic(fmt.Sprintf("%s is nil", name))
		}
	}
}

package assert

import "fmt"

func AssertEqual[T comparable](expected T, received T, msg string) {
	if expected != received {
		msg = fmt.Sprintf("%s. Expected: %v. Receive: %v", msg, expected, received)
		panic(msg)
	}
}

func AssertTrue(expr bool, msg string) {
	if !expr {
		panic(msg)
	}
}

func AssertNil(target any, name string) {
	if target != nil {
		msg := fmt.Sprintf("%s is not nil", name)
		panic(msg)
	}
}

func AssertNotNil(target any, name string) {
	if target == nil {
		msg := fmt.Sprintf("%s is nil!", name)
		panic(msg)
	}
}

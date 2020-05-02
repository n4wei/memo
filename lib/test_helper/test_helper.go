package test_helper

import (
	"reflect"
	"testing"
)

func AssertEqual(t *testing.T, actual, expected interface{}) {
	if actual != expected {
		printError(t, actual, expected)
	}
}

func AssertDeepEqual(t *testing.T, actual, expected interface{}) {
	if !reflect.DeepEqual(actual, expected) {
		printError(t, actual, expected)
	}
}

func AssertError(t *testing.T, err error) {
	if err == nil {
		printError(t, nil, "an error")
	}
}

func AssertNoError(t *testing.T, err error) {
	if err != nil {
		printError(t, err, nil)
	}
}

func printError(t *testing.T, actual, expected interface{}) {
	t.Errorf("\nActual:\n%v\nExpected:\n%v\n", actual, expected)
}

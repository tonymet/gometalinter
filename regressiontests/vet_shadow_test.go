package regressiontests

import (
	"runtime"
	"strings"
	"testing"
)

func TestVetShadow(t *testing.T) {
	t.Parallel()
	source := `package test

type MyStruct struct {}
func test(mystructs []*MyStruct) *MyStruct {
	var foo *MyStruct
	for _, mystruct := range mystructs {
		foo := mystruct
		_ = foo
	}
	return foo
}
`
	expected := Issues{
		{Linter: "vetshadow", Severity: "warning", Path: "test.go", Line: 7, Col: 3, Message: `declaration of "foo" shadows declaration at line 5`},
	}

	if version := runtime.Version(); strings.HasPrefix(version, "go1.8") {
		expected = Issues{
			{Linter: "vetshadow", Severity: "warning", Path: "test.go", Line: 7, Col: 3, Message: "foo declared and not used"},
		}
	}

	ExpectIssues(t, "vetshadow", source, expected)
}

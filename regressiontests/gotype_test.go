package regressiontests

import (
	"fmt"
	"testing"

	"gotest.tools/v3/fs"
	"github.com/stretchr/testify/assert"
)

func TestGoType(t *testing.T) {
	t.Parallel()

	dir := fs.NewDir(t, "test-gotype",
		fs.WithFile("file.go", goTypeFile("root")),
		fs.WithDir("sub",
			fs.WithFile("file.go", goTypeFile("sub"))),
		fs.WithDir("excluded",
			fs.WithFile("file.go", goTypeFile("excluded"))))
	defer dir.Remove()

	expected := Issues{
		{Linter: "gotype", Severity: "error", Path: "file.go", Line: 4, Col: 6, Message: "declared and not used: foo"},
		{Linter: "gotype", Severity: "error", Path: "sub/file.go", Line: 4, Col: 6, Message: "declared and not used: foo"},
	}
	actual := RunLinter(t, "gotype", dir.Path(), "--skip=excluded")
	assert.Equal(t, expected, actual)
}

func TestGoTypeWithMultiPackageDirectoryTest(t *testing.T) {
	t.Parallel()

	dir := fs.NewDir(t, "test-gotype",
		fs.WithFile("file.go", goTypeFile("root")),
		fs.WithFile("file_test.go", goTypeFile("root_test")))
	defer dir.Remove()

	expected := Issues{
		{Linter: "gotype", Severity: "error", Path: "file.go", Line: 4, Col: 6, Message: "declared and not used: foo"},
		{Linter: "gotypex", Severity: "error", Path: "file_test.go", Line: 4, Col: 6, Message: "declared and not used: foo"},
	}
	actual := RunLinter(t, "gotype", dir.Path())
	actual = append(actual, RunLinter(t, "gotypex", dir.Path())...)
	assert.Equal(t, expected, actual)
}


func goTypeFile(pkg string) string {
	return fmt.Sprintf(`package %s

func badFunction() {
	var foo string
}
	`, pkg)
}

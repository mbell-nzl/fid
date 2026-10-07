package main

import (
	"go/ast"
	"go/parser"
	"go/token"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

// R9 is an adoption audit: raw helpers remain only for the two explicitly
// excluded diagnostics. Process tests prove the six newly adopted paths.
func TestP24HumanFailureCallerAudit(t *testing.T) {
	t.Parallel()
	files, err := filepath.Glob("*.go")
	if err != nil {
		t.Fatal(err)
	}
	for _, path := range files {
		if strings.HasSuffix(path, "_test.go") {
			continue
		}
		file, err := parser.ParseFile(token.NewFileSet(), path, nil, 0)
		if err != nil {
			t.Fatal(err)
		}
		ast.Inspect(file, func(node ast.Node) bool {
			call, ok := node.(*ast.CallExpr)
			if !ok {
				return true
			}
			if path == "root.go" {
				for _, arg := range call.Args {
					if value, ok := arg.(*ast.BasicLit); ok && value.Kind == token.STRING {
						text, err := strconv.Unquote(value.Value)
						if err != nil {
							t.Fatal(err)
						}
						if strings.HasPrefix(text, "error: ") {
							t.Errorf("root-flag failure bypasses shared renderer: %s", path)
						}
					}
				}
			}
			fn, ok := call.Fun.(*ast.Ident)
			if !ok {
				return true
			}
			switch fn.Name {
			case "renderIssueError":
				if path != "journal.go" && path != "sql.go" {
					t.Errorf("stable caller bypasses shared failure renderer: %s", path)
				}
			case "renderIssueReadError":
				if len(call.Args) > 1 {
					if value, ok := call.Args[1].(*ast.Ident); ok && value.Name == "true" {
						t.Errorf("caller forces JSON failure mode: %s", path)
					}
				}
			}
			return true
		})
	}
}

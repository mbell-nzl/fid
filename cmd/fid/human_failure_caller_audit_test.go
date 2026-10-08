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
		if path == "root.go" {
			if count := p24RawRootFailures(file, false); count != 0 {
				t.Errorf("root-flag failure bypasses shared renderer: %s (%d)", path, count)
			}
		}
		ast.Inspect(file, func(node ast.Node) bool {
			call, ok := node.(*ast.CallExpr)
			if !ok {
				return true
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

// Follow the AST branch, rather than allowing raw text by filename or function.
// Only the positive diagnosticLeaf branch permits the retained raw refusal.
func p24RawRootFailures(node ast.Node, diagnostic bool) int {
	count := 0
	ast.Inspect(node, func(node ast.Node) bool {
		if branch, ok := node.(*ast.IfStmt); ok {
			permitted := false
			if call, ok := branch.Cond.(*ast.CallExpr); ok {
				if fn, ok := call.Fun.(*ast.Ident); ok && fn.Name == "diagnosticLeaf" {
					permitted = true
				}
			}
			count += p24RawRootFailures(branch.Cond, diagnostic)
			count += p24RawRootFailures(branch.Body, diagnostic || permitted)
			if branch.Else != nil {
				count += p24RawRootFailures(branch.Else, diagnostic)
			}
			return false
		}
		if literal, ok := node.(*ast.BasicLit); ok && literal.Kind == token.STRING && !diagnostic {
			text, err := strconv.Unquote(literal.Value)
			if err == nil && strings.HasPrefix(text, "error: ") {
				count++
			}
		}
		return true
	})
	return count
}

func TestP24RootFailureAuditNegativeControl(t *testing.T) {
	t.Parallel()
	for _, row := range []struct {
		source string
		want   int
	}{
		{`package main; func reject(c any) { if diagnosticLeaf(c) { println("error: diagnostic") } }`, 0},
		{`package main; func reject(c any) { if diagnosticLeaf(c) { println("error: diagnostic") } else { println("error: stable") } }`, 1},
		{`package main; func reject(c any) { println("error: stable") }`, 1},
		{`package main; func reject(c any) { if !diagnosticLeaf(c) { println("error: stable") } }`, 1},
	} {
		file, err := parser.ParseFile(token.NewFileSet(), "root.go", row.source, 0)
		if err != nil {
			t.Fatal(err)
		}
		if got := p24RawRootFailures(file, false); got != row.want {
			t.Fatalf("raw failure count=%d want=%d: %s", got, row.want, row.source)
		}
	}
}

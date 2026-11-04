package config

import (
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

// The ledger is moving to the Git common directory (RESET task 1.1d), so the
// only code that may join the ledger directory onto a root is LedgerDir. That
// holds for tests too, so they follow the ledger when it moves. The config
// package's own tests are exempt: they pin where LedgerDir points. A path that
// starts with ".fid" is a display string, not a location, and is left alone;
// one joined onto a root is a location.
func TestLedgerPathsComeOnlyFromLedgerDir(t *testing.T) {
	t.Parallel()
	var strays []string
	for _, tree := range []string{"../../cmd", "../../internal"} {
		err := filepath.WalkDir(tree, func(path string, entry fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if entry.IsDir() && entry.Name() == "testdata" {
				return filepath.SkipDir
			}
			if entry.IsDir() || !strings.HasSuffix(path, ".go") {
				return nil
			}
			found, err := strayLedgerJoins(path)
			strays = append(strays, found...)
			return err
		})
		if err != nil {
			t.Fatalf("walk %s: %v", tree, err)
		}
	}
	if len(strays) != 0 {
		t.Errorf("%d ledger paths are built without LedgerDir:\n%s", len(strays), strings.Join(strays, "\n"))
	}
}

func strayLedgerJoins(path string) ([]string, error) {
	files := token.NewFileSet()
	file, err := parser.ParseFile(files, path, nil, 0)
	if err != nil {
		return nil, err
	}
	if file.Name.Name == "config" && strings.HasSuffix(path, "_test.go") {
		return nil, nil
	}
	var strays []string
	ast.Inspect(file, func(node ast.Node) bool {
		if decl, ok := node.(*ast.FuncDecl); ok && file.Name.Name == "config" && decl.Recv == nil && decl.Name.Name == "LedgerDir" {
			return false
		}
		call, ok := node.(*ast.CallExpr)
		if !ok || !isFilepathJoin(call) || len(call.Args) < 2 {
			return true
		}
		for _, argument := range call.Args[1:] {
			if namesLedgerDirectory(argument) {
				strays = append(strays, files.Position(argument.Pos()).String())
			}
		}
		return true
	})
	return strays, nil
}

func isFilepathJoin(call *ast.CallExpr) bool {
	selector, ok := call.Fun.(*ast.SelectorExpr)
	if !ok || selector.Sel.Name != "Join" {
		return false
	}
	pkg, ok := selector.X.(*ast.Ident)
	return ok && pkg.Name == "filepath"
}

func namesLedgerDirectory(expression ast.Expr) bool {
	switch value := expression.(type) {
	case *ast.BasicLit:
		text, err := strconv.Unquote(value.Value)
		return value.Kind == token.STRING && err == nil && (text == ".fid" || strings.HasPrefix(text, ".fid/"))
	case *ast.SelectorExpr:
		return value.Sel.Name == "WorkspaceDirName"
	case *ast.Ident:
		return value.Name == "WorkspaceDirName"
	}
	return false
}

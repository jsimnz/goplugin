package subcommands

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/printer"
	"go/token"
	"os"
	"path"
)

func getNodeString(n ast.Node, fset *token.FileSet) string {
	var buf bytes.Buffer
	err := printer.Fprint(&buf, fset, n)
	if err != nil {
		return ""
	}
	return string(buf.Bytes())
}

func Abs(name string) (string, error) {
	if path.IsAbs(name) {
		return name, nil
	}
	wd, err := os.Getwd()
	return path.Join(wd, name), err
}

func getIdentType(n *ast.Ident) string {
	return (((n.Obj.Decl).(*ast.ValueSpec).Values[0]).(*ast.CompositeLit).Type).(*ast.Ident).Name
}

func stringer(i interface{}) string {
	return fmt.Sprintf("%#v", i)
}

func typeToString(typ ast.Expr) string {
	switch s := (typ).(type) {
	case *ast.Ident:
		return s.Name
	case *ast.StarExpr:
		return "*" + (s.X).(*ast.Ident).Name
	}
	return ""
}

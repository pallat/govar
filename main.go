package main

import (
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
)

type visitor struct {
	fn string
}

func (v *visitor) Visit(node ast.Node) (w ast.Visitor) {
	switch t := node.(type) {
	case *ast.GenDecl:
		for _, spec := range t.Specs {
			switch spec.(type) {
			case *ast.ValueSpec:
				for _, value := range spec.(*ast.ValueSpec).Values {
					switch value.(type) {
					case *ast.CallExpr:
						switch value.(*ast.CallExpr).Fun.(type) {
						case *ast.SelectorExpr:
							switch value.(*ast.CallExpr).Fun.(*ast.SelectorExpr).X.(type) {
							case *ast.Ident:
								if value.(*ast.CallExpr).Fun.(*ast.SelectorExpr).X.(*ast.Ident).Name == pack {
									fmt.Println(v.fn, " : ", value.(*ast.CallExpr).Args[0].(*ast.BasicLit).Value)
								}
							}
						}
					}
				}
			}
		}
	}

	return v
}

func walk(path string, info os.FileInfo, err error) error {
	p := strings.Split(path, ".")
	if p[len(p)-1] != "go" {
		return nil
	}
	find(path)
	return nil
}

var root string
var pack string

func init() {
	flag.StringVar(&root, "root", ".", "-root=/")
	flag.StringVar(&pack, "pack", "conf", "-pack=conf")
}

func main() {
	flag.Parse()
	filepath.Walk(root, walk)
}

func find(fn string) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, fn, nil, parser.AllErrors)
	if err != nil {
		fmt.Println(err)
		return
	}

	// ast.Print(fset, f)

	ast.Walk(&visitor{fn: fn}, f)
}

package main

import (
	"go/ast"
	"go/parser"
	"go/token"
)

func createAsts(files map[string]string) (map[string]*ast.File, error) {
	set := token.NewFileSet()
	asts := make(map[string]*ast.File)

	for name, source := range files {
		node, err := parser.ParseFile(set, name, source, parser.AllErrors)
		if err != nil {
			return nil, err
		}
		asts[name] = node
	}

	return asts, nil
}

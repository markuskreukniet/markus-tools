package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"reflect"
)

func parseFileDefaultMode(set *token.FileSet, name, source string) (*ast.File, error) {
	return parser.ParseFile(set, name, source, 0)
}

func createAsts(files map[string]string) (map[string]*ast.File, error) {
	set := token.NewFileSet()
	asts := make(map[string]*ast.File)

	for name, source := range files {
		node, err := parseFileDefaultMode(set, name, source)
		if err != nil {
			return nil, err
		}
		asts[name] = node
	}

	return asts, nil
}

// An Obj represents a reference to a declared symbol in the code.
// Symbols include variables, parameters, functions, constants, types, and labels.
func areFunctionsStructurallyEqual(declI, declJ *ast.FuncDecl) bool {
	declI.Name = &ast.Ident{Name: ""}
	declJ.Name = &ast.Ident{Name: ""}

	clearParameterNamesAndObjects := func(decl *ast.FuncDecl) {
		if decl.Type != nil && decl.Type.Params != nil {
			for _, field := range decl.Type.Params.List {
				for _, name := range field.Names {
					name.Name = ""
					name.Obj = nil
				}
			}
		}
	}

	clearFunctionBodyBoundIdentifiers := func(decl *ast.FuncDecl) {
		ast.Inspect(decl, func(n ast.Node) bool {
			if ident, ok := n.(*ast.Ident); ok && ident.Obj != nil {
				ident.Name = ""
				ident.Obj = nil
			}
			return true
		})
	}

	clearParameterNamesAndObjects(declI)
	clearParameterNamesAndObjects(declJ)

	clearFunctionBodyBoundIdentifiers(declI)
	clearFunctionBodyBoundIdentifiers(declJ)

	declI = clearPositions(declI).(*ast.FuncDecl)
	declJ = clearPositions(declJ).(*ast.FuncDecl)

	return reflect.DeepEqual(declI, declJ)
}

func clearPositions(n ast.Node) ast.Node {
	// token.Pos(0) is a Pos value with underlying integer value 0, representing an undefined position.
	undefinedPos := token.Pos(0)
	posType := reflect.TypeOf(undefinedPos)            // runtime type of token.Pos
	undefinedPosValue := reflect.ValueOf(undefinedPos) // runtime value of undefined token.Pos

	ast.Inspect(n, func(n ast.Node) bool {
		if n == nil {
			return false // skip this node and its children
		}

		// access the struct behind the interface and pointer. For example, elem could represent a full *ast.FuncDecl.
		elem := reflect.ValueOf(n).Elem()
		for i := 0; i < elem.NumField(); i++ {
			field := elem.Field(i)

			// CanSet is required to avoid panic when trying to set unexported (lowercase) fields.
			if field.Type() == posType && field.CanSet() {
				field.Set(undefinedPosValue)
			}
		}
		return true
	})

	return n
}

// TODO: below here

type codeLocation struct {
	StartLine   int
	StartColumn int
	EndLine     int
	EndColumn   int
}

type duplicateCodeParts [][]codeLocation

type duplicateStatements struct {
	originalStatementIndex int
	originalStatement      ast.Stmt
	duplicateStatements    []ast.Stmt
}

func findStructuralDuplicateFunctionBodyParts(set *token.FileSet, declI, declJ *ast.FuncDecl) (duplicateCodeParts, error) {
	if len(declI.Body.List) == 0 || len(declJ.Body.List) == 0 {
		return nil, nil
	}

	var parts duplicateCodeParts
	var duplicateStatementsList []duplicateStatements

	for i, stmtI := range declI.Body.List {
		var foundStatements []ast.Stmt

		for _, stmtJ := range declJ.Body.List {
			cloneI, err := cloneStatement(stmtI)
			if err != nil {
				return nil, err
			}

			cloneJ, err := cloneStatement(stmtJ)
			if err != nil {
				return nil, err
			}

			cloneI = clearPositions(cloneI).(ast.Stmt) // TODO: possible to skip this casting?
			cloneJ = clearPositions(cloneJ).(ast.Stmt)

			if reflect.DeepEqual(cloneI, cloneJ) {
				foundStatements = append(foundStatements, stmtJ)
			}
		}

		if len(foundStatements) > 0 {
			duplicateStatementsList = append(duplicateStatementsList, duplicateStatements{
				originalStatementIndex: i,
				originalStatement:      stmtI,
				duplicateStatements:    foundStatements,
			})
			foundStatements = foundStatements[:0] // TODO: search for " = nill", maybe use [:0] more?
		}
	}

	var group []duplicateStatements
	lastIndex := -2
	for _, item := range duplicateStatementsList {
		length := len(group)

		if lastIndex == item.originalStatementIndex-1 {
			group = append(group, item)
		} else if length > 0 {
			start := set.Position(group[0].originalStatement.Pos())
			end := set.Position(group[length-1].originalStatement.End())

			locationI := codeLocation{
				StartLine:   start.Line,
				StartColumn: start.Column,
				EndLine:     end.Line,
				EndColumn:   end.Column,
			}

			start = set.Position(group[0].duplicateStatements[0].Pos())
			end = set.Position(group[length-1].duplicateStatements[0].End())

			locationJ := codeLocation{
				StartLine:   start.Line,
				StartColumn: start.Column,
				EndLine:     end.Line,
				EndColumn:   end.Column,
			}

			parts = append(parts, []codeLocation{locationI, locationJ})

			group = group[:0]
		} else {
			group = []duplicateStatements{item}
		}
	}

	return parts, nil
}

func cloneStatement(stmt ast.Stmt) (ast.Stmt, error) {
	var buffer bytes.Buffer
	if err := format.Node(&buffer, token.NewFileSet(), stmt); err != nil {
		return nil, err
	}

	file, err := parseFileDefaultMode(
		token.NewFileSet(), "", fmt.Sprintf("package main\nfunc dummy() {\n%s\n}", buffer.String()),
	)
	if err != nil {
		return nil, err
	}

	return file.Decls[0].(*ast.FuncDecl).Body.List[0], nil
}

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

// func printNode(label string, node ast.Node) {
// 	println("======", label)
// 	_ = printer.Fprint(os.Stdout, token.NewFileSet(), node)
// 	println("\n======")
// }

func parseFileDefaultMode(set *token.FileSet, name, source string) (*ast.File, error) {
	return parser.ParseFile(set, name, source, 0)
}

func createAsts(set *token.FileSet, files map[string]string) (map[string]*ast.File, error) {
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

// TODO: check if other acronyms on other places are in all caps.
// TODO: use this syntax on other places such as must..?
func normalizeASTNodes(nodes ...ast.Node) {
	for _, node := range nodes {
		clearIdentifiers(node)
		clearPositions(node)
	}
}

func clearIdentifiers(n ast.Node) {
	ast.Inspect(n, func(n ast.Node) bool {
		// if expr, ok := n.(*ast.CallExpr); ok {
		// 	if ident, ok := expr.Fun.(*ast.Ident); ok {
		// 		ident.Obj = nil
		// 		return true
		// 	}
		// }

		if ident, ok := n.(*ast.Ident); ok {
			ident.Name = ""
			ident.Obj = nil
		}
		return true
	})
}

func clearPositions(n ast.Node) {
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

	clearParameterNamesAndObjects(declI)
	clearParameterNamesAndObjects(declJ)

	normalizeASTNodes(declI, declJ)

	return reflect.DeepEqual(declI, declJ)
}

type codeLocation struct {
	startLine   int
	startColumn int
	endLine     int
	endColumn   int
}

func createCodeLocation(start, end token.Position) codeLocation {
	return codeLocation{
		startLine:   start.Line,
		startColumn: start.Column,
		endLine:     end.Line,
		endColumn:   end.Column,
	}
}

type duplicateStatement struct {
	originalIndex  int
	duplicateIndex int
	original       ast.Stmt
	duplicate      ast.Stmt
}

func findStructuralDuplicateFunctionParts(
	set *token.FileSet, declI, declJ *ast.FuncDecl,
) (bool, [][]codeLocation, error) {
	var duplicateStatements []duplicateStatement
	areBodiesTheSame := true

	for i, stmtI := range declI.Body.List {
		for j, stmtJ := range declJ.Body.List {
			cloneI, err := cloneStatement(stmtI)
			if err != nil {
				return false, nil, err
			}

			cloneJ, err := cloneStatement(stmtJ)
			if err != nil {
				return false, nil, err
			}

			normalizeASTNodes(cloneI, cloneJ)

			if reflect.DeepEqual(cloneI, cloneJ) {
				duplicateStatements = append(duplicateStatements, duplicateStatement{
					originalIndex:  i,
					duplicateIndex: j,
					original:       stmtI,
					duplicate:      stmtJ,
				})
			} else {
				areBodiesTheSame = false // TODO: not efficient
			}
		}
	}

	var parts [][]codeLocation
	var duplicateStatementGroup []duplicateStatement
	lastOriginalIndex, lastDuplicateIndex := -1, -1

	appendPart := func(length int) {
		parts = append(parts, []codeLocation{
			createCodeLocation(
				set.Position(duplicateStatementGroup[0].original.Pos()),
				set.Position(duplicateStatementGroup[length-1].original.Pos()),
			),
			createCodeLocation(
				set.Position(duplicateStatementGroup[0].duplicate.Pos()),
				set.Position(duplicateStatementGroup[length-1].duplicate.Pos()),
			),
		})
	}

	for _, statement := range duplicateStatements {
		if lastOriginalIndex+1 == statement.originalIndex && lastDuplicateIndex+1 == statement.duplicateIndex {
			duplicateStatementGroup = append(duplicateStatementGroup, statement)
		} else {
			appendPart(len(duplicateStatementGroup))
			duplicateStatementGroup = duplicateStatementGroup[:0]
		}
	}

	length := len(duplicateStatementGroup)
	if length > 0 {
		appendPart(length)
	}

	return areBodiesTheSame, parts, nil
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

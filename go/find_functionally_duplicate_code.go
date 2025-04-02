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

func createASTs(set *token.FileSet, files map[string]string) (map[string]*ast.File, error) {
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

// TODO: use this param syntax on other places such as must..?
func normalizeASTNodes(nodes ...ast.Node) {
	for _, node := range nodes {
		normalizeIdentifiers(node)
		clearPositions(node)
	}
}

// TODO: add comments
func normalizeIdentifiers(n ast.Node) {
	preserved := make(map[*ast.Ident]struct{})

	ast.Inspect(n, func(n ast.Node) bool {
		switch node := n.(type) {
		case *ast.CallExpr:
			switch expr := node.Fun.(type) {
			case *ast.Ident:
				preserved[expr] = struct{}{} // Preserve direct function calls: do()
			case *ast.SelectorExpr:
				// Preserve method receiver: obj in obj.Do()
				if ident, ok := expr.X.(*ast.Ident); ok {
					preserved[ident] = struct{}{}
				}
				// Preserve selector: Do in obj.Do() or Println in fmt.Println()
				if expr.Sel != nil {
					preserved[expr.Sel] = struct{}{}
				}
			}
		case *ast.Ident:
			if _, ok := preserved[node]; ok {
				return true
			}
			node.Name = ""
			node.Obj = nil
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

type functionFile struct {
	fileName  string
	functions []*ast.FuncDecl
}

func findAllFunctions(asts map[string]*ast.File) []functionFile {
	var files []functionFile

	for fileName, file := range asts {
		var functions []*ast.FuncDecl
		for _, decl := range file.Decls {
			if funcDecl, ok := decl.(*ast.FuncDecl); ok {
				functions = append(functions, funcDecl)
			}
		}
		files = append(files, functionFile{
			fileName:  fileName,
			functions: functions,
		})
	}

	return files
}

// func findFunctionallyEqualFunctions(asts map[string]*ast.File) {
// 	files := findAllFunctions(asts)

// 	for i, fileI := range files {
// 		for _, fileJ := range files[i+1:] {
// 			for k, functionI := range fileI.functions {
// 				for l, functionJ := range fileJ.functions {
// 					if k == l {
// 						continue
// 					}

// 					// TODO: should not do anonymizeFuncDeclHeader and normalizeASTNodes every time?
// 					if areFunctionsFunctionallyEqual(functionI, functionJ) {
// 						//
// 					}
// 				}
// 			}
// 		}
// 	}
// }

// An Obj represents a reference to a declared symbol in the code.
// Symbols include variables, parameters, functions, constants, types, and labels.
func areFunctionsFunctionallyEqual(declI, declJ *ast.FuncDecl) bool {
	anonymizeFuncDeclHeader := func(decl *ast.FuncDecl) {
		decl.Name = &ast.Ident{Name: ""}

		if decl.Type != nil && decl.Type.Params != nil {
			for _, field := range decl.Type.Params.List {
				for _, name := range field.Names {
					name.Name = ""
					name.Obj = nil
				}
			}
		}
	}

	anonymizeFuncDeclHeader(declI)
	anonymizeFuncDeclHeader(declJ)

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

func findDuplicateStatements(statementsI, statementsJ []ast.Stmt) ([]duplicateStatement, error) {
	var duplicates []duplicateStatement

	for i, stmtI := range statementsI {
		cloneI, err := cloneStatement(stmtI)
		if err != nil {
			return nil, err
		}

		for j, stmtJ := range statementsJ {
			cloneJ, err := cloneStatement(stmtJ)
			if err != nil {
				return nil, err
			}

			normalizeASTNodes(cloneI, cloneJ)

			if reflect.DeepEqual(cloneI, cloneJ) {
				duplicates = append(duplicates, duplicateStatement{
					originalIndex:  i,
					duplicateIndex: j,
					original:       stmtI,
					duplicate:      stmtJ,
				})
			}
		}
	}

	return duplicates, nil
}

func findDuplicateCodeRegions(set *token.FileSet, declI, declJ *ast.FuncDecl) ([][]codeLocation, error) {
	duplicateStatements, err := findDuplicateStatements(declI.Body.List, declJ.Body.List)
	if err != nil {
		return nil, err
	}

	var duplicateCodeRegions [][]codeLocation
	var duplicateStatementGroup []duplicateStatement
	lastOriginalIndex, lastDuplicateIndex := -1, -1

	appendRegion := func(length int) {
		duplicateCodeRegions = append(duplicateCodeRegions, []codeLocation{
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
			appendRegion(len(duplicateStatementGroup))
			duplicateStatementGroup = duplicateStatementGroup[:0]
		}
		lastOriginalIndex = statement.originalIndex
		lastDuplicateIndex = statement.duplicateIndex
	}

	length := len(duplicateStatementGroup)
	if length > 0 {
		appendRegion(length)
	}

	return duplicateCodeRegions, nil
}

// TODO: function still needed?
func cloneStatement(stmt ast.Stmt) (ast.Stmt, error) {
	var buffer bytes.Buffer
	set := token.NewFileSet()

	if err := format.Node(&buffer, set, stmt); err != nil {
		return nil, err
	}

	file, err := parseFileDefaultMode(set, "", fmt.Sprintf("package main\nfunc dummy() {\n%s\n}", buffer.String()))
	if err != nil {
		return nil, err
	}

	return file.Decls[0].(*ast.FuncDecl).Body.List[0], nil
}

// TODO: WIP
// func cloneFuncDecl(decl *ast.FuncDecl) (*ast.FuncDecl, error) {
// 	var buffer bytes.Buffer
// 	set := token.NewFileSet()

// 	if err := format.Node(&buffer, set, decl); err != nil {
// 		return nil, err
// 	}

// 	file, err := parser.ParseFile(set, "", fmt.Sprintf("package main\n%s", buffer.String()), 0)
// 	if err != nil {
// 		return nil, fmt.Errorf("parsing failed: %w", err)
// 	}

// 	// Extract the first declaration, assuming it's a function
// 	if len(file.Decls) == 0 {
// 		return nil, fmt.Errorf("no declarations found")
// 	}

// 	cloned, ok := file.Decls[0].(*ast.FuncDecl)
// 	if !ok {
// 		return nil, fmt.Errorf("decl is not a FuncDecl")
// 	}

// 	return cloned, nil
// }

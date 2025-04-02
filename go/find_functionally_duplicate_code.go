package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"reflect"
	"strings"

	"github.com/markuskreukniet/markus-tools/go/utils"
)

// func printNode(label string, node ast.Node) {
// 	println("=====", label)
// 	_ = printer.Fprint(os.Stdout, token.NewFileSet(), node)
// 	println("\n=====")
// }

// An Obj represents a reference to a declared symbol in the code.
// Symbols include variables, parameters, functions, constants, types, and labels.

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

func normalizeIdentifiers(n ast.Node) {
	// ast.Inspect performs a depth-first traversal, so the same identifier node may be visited multiple times.
	// We preserve identifiers that are part of function or method calls to avoid clearing them during traversal.
	// When these identifiers are visited again as standalone *ast.Ident nodes, we skip them.
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

type normalizedFunction struct {
	function           *ast.FuncDecl
	normalizedFunction *ast.FuncDecl
	parameterTypeCount map[string]int
	resultTypeCount    map[string]int
}

type normalizedFunctionFile struct {
	fileName  string
	functions []normalizedFunction
}

func extractNormalizedFunctions(set *token.FileSet, asts map[string]*ast.File) ([]normalizedFunctionFile, error) {
	var files []normalizedFunctionFile

	for fileName, file := range asts {
		var functions []normalizedFunction
		for _, decl := range file.Decls {
			if funcDecl, ok := decl.(*ast.FuncDecl); ok {
				clone, err := cloneFuncDecl(funcDecl)
				if err != nil {
					return nil, err
				}

				parameterTypeCount, resultTypeCount, err := normalizeAndCountFunctionTypes(set, clone)
				if err != nil {
					return nil, err
				}

				functions = append(functions, normalizedFunction{
					function:           funcDecl,
					normalizedFunction: clone,
					parameterTypeCount: parameterTypeCount,
					resultTypeCount:    resultTypeCount,
				})
			}
		}
		files = append(files, normalizedFunctionFile{
			fileName:  fileName,
			functions: functions,
		})
	}

	return files, nil
}

func normalizeAndCountFunctionTypes(set *token.FileSet, decl *ast.FuncDecl) (map[string]int, map[string]int, error) {
	countAndClearFieldListTypes := func(fields []*ast.Field, typeCount *map[string]int, list **ast.FieldList) error {
		for _, field := range fields {
			var buffer bytes.Buffer
			if err := format.Node(&buffer, set, field.Type); err != nil {
				return err
			}
			typeString := buffer.String()
			length := len(field.Names)
			if length == 0 { // anonymous parameter or result (no name specified)
				(*typeCount)[typeString]++
			} else {
				(*typeCount)[typeString] += length
			}
		}

		*list = nil

		return nil
	}

	parameterTypeCount, resultTypeCount := make(map[string]int), make(map[string]int)

	decl.Name = &ast.Ident{Name: ""}

	if decl.Type != nil {
		if decl.Type.Params != nil {
			if err := countAndClearFieldListTypes(
				decl.Type.Params.List, &parameterTypeCount, &decl.Type.Params,
			); err != nil {
				return nil, nil, err
			}
		}
		if decl.Type.Results != nil {
			if err := countAndClearFieldListTypes(
				decl.Type.Results.List, &parameterTypeCount, &decl.Type.Results,
			); err != nil {
				return nil, nil, err
			}
		}
	}

	normalizeASTNodes(decl)

	return parameterTypeCount, resultTypeCount, nil
}

func areTypeCountsEqual(countsI, countsJ map[string]int) bool {
	if len(countsI) != len(countsJ) {
		return false
	}

	for key, valI := range countsI {
		if valJ, ok := countsJ[key]; !ok || valJ != valI {
			return false
		}
	}

	return true
}

func findFunctionallyEqualFunctions(set *token.FileSet, asts map[string]*ast.File) (string, error) {
	files, err := extractNormalizedFunctions(set, asts)
	if err != nil {
		return "", err
	}

	var buffer bytes.Buffer

	for i, fileI := range files {
		for j, fileJ := range files[i:] {
			for k, functionI := range fileI.functions {
				for l, functionJ := range fileJ.functions {
					if j == 0 && k == l {
						continue
					}

					if areTypeCountsEqual(functionI.parameterTypeCount, functionJ.parameterTypeCount) &&
						areTypeCountsEqual(functionI.resultTypeCount, functionJ.resultTypeCount) &&
						reflect.DeepEqual(functionI.normalizedFunction, functionJ.normalizedFunction) {
						if _, err := buffer.WriteString(
							codePositionsToText(
								set,
								fileI.fileName, fileJ.fileName,
								functionI.function.Pos(), functionI.function.End(),
								functionJ.function.Pos(), functionJ.function.End(),
							),
						); err != nil {
							return "", err
						}
						if _, err := utils.WriteTwoNewlineStrings(&buffer); err != nil {
							return "", err
						}
					}
				}
			}
		}
	}

	length := buffer.Len()
	if length >= 2 {
		buffer.Truncate(length - 2)
	}

	return buffer.String(), nil
}

func codePositionsToText(
	set *token.FileSet,
	fileNameOriginal, fileNameDuplicate string,
	startOriginal, endOriginal, startDuplicate, endDuplicate token.Pos,
) string {
	text := `file name: %s
Start line %d and column %d.
End line %d and column %d.
file name: %s
Start line %d and column %d.
End line %d and column %d.`

	return fmt.Sprintf(
		text,
		fileNameOriginal,
		set.Position(startOriginal).Line, set.Position(startOriginal).Column,
		set.Position(endOriginal).Line, set.Position(endOriginal).Column,
		fileNameDuplicate,
		set.Position(startDuplicate).Line, set.Position(startDuplicate).Column,
		set.Position(endDuplicate).Line, set.Position(endDuplicate).Column,
	)
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

func findDuplicateCodeRegions(set *token.FileSet, declI, declJ *ast.FuncDecl) (string, error) {
	duplicateStatements, err := findDuplicateStatements(declI.Body.List, declJ.Body.List)
	if err != nil {
		return "", err
	}

	var builder strings.Builder
	var duplicateStatementGroup []duplicateStatement
	lastOriginalIndex, lastDuplicateIndex := -1, -1

	writeCodePositions := func(length int) error {
		if _, err := builder.WriteString(
			codePositionsToText(
				set,
				"", "",
				duplicateStatementGroup[0].original.Pos(), duplicateStatementGroup[length-1].original.End(),
				duplicateStatementGroup[0].duplicate.Pos(), duplicateStatementGroup[length-1].duplicate.End(),
			),
		); err != nil {
			return err
		}

		return nil
	}

	for _, statement := range duplicateStatements {
		if lastOriginalIndex+1 == statement.originalIndex && lastDuplicateIndex+1 == statement.duplicateIndex {
			duplicateStatementGroup = append(duplicateStatementGroup, statement)
		} else {
			writeCodePositions(len(duplicateStatementGroup))
			utils.WriteTwoNewlineStrings(&builder)
			duplicateStatementGroup = duplicateStatementGroup[:0]
		}
		lastOriginalIndex = statement.originalIndex
		lastDuplicateIndex = statement.duplicateIndex
	}

	length := len(duplicateStatementGroup)
	if length > 0 {
		writeCodePositions(length)
	}

	return builder.String(), nil
}

func cloneStatement(stmt ast.Stmt) (ast.Stmt, error) {
	generateGoSource := func(buffer bytes.Buffer) string {
		return fmt.Sprintf("package main\nfunc dummy() {\n%s\n}", buffer.String())
	}

	retrieveTargetNode := func(funcDecl *ast.FuncDecl) ast.Node {
		return funcDecl.Body.List[0]
	}

	return cloneNodeViaFuncDecl[ast.Stmt](stmt, generateGoSource, retrieveTargetNode)
}

func cloneFuncDecl(decl *ast.FuncDecl) (*ast.FuncDecl, error) {
	generateGoSource := func(buffer bytes.Buffer) string {
		return fmt.Sprintf("package main\n%s", buffer.String())
	}

	retrieveTargetNode := func(funcDecl *ast.FuncDecl) ast.Node {
		return funcDecl
	}

	return cloneNodeViaFuncDecl[*ast.FuncDecl](decl, generateGoSource, retrieveTargetNode)
}

func cloneNodeViaFuncDecl[T ast.Node](
	n ast.Node, handlerI func(bytes.Buffer) string, handlerJ func(*ast.FuncDecl) ast.Node,
) (T, error) {
	var zero T
	var buffer bytes.Buffer
	set := token.NewFileSet()

	if err := format.Node(&buffer, set, n); err != nil {
		return zero, err
	}

	file, err := parseFileDefaultMode(set, "", handlerI(buffer))
	if err != nil {
		return zero, err
	}

	return handlerJ(file.Decls[0].(*ast.FuncDecl)).(T), nil
}

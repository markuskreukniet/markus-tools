package main

import (
	"go/ast"
	"go/token"
	"reflect"
	"testing"

	"github.com/markuskreukniet/markus-tools/go/utils"
)

func tMustCreateAsts(t *testing.T, files map[string]string) map[string]*ast.File {
	result, err := createAsts(files)
	return utils.TMust(t, result, err)
}

// An Obj represents a reference to a declared symbol in the code.
// Symbols include variables, parameters, functions, constants, types, and labels.
func areFunctionsStructurallyEqual(a, b *ast.FuncDecl) bool {
	a.Name = &ast.Ident{Name: ""}
	b.Name = &ast.Ident{Name: ""}

	clearParameterNamesAndObjects := func(d *ast.FuncDecl) {
		if d.Type != nil && d.Type.Params != nil {
			for _, field := range d.Type.Params.List {
				for _, name := range field.Names {
					name.Name = ""
					name.Obj = nil
				}
			}
		}
	}

	clearFunctionBodyBoundIdentifiers := func(d *ast.FuncDecl) {
		ast.Inspect(d, func(n ast.Node) bool {
			if ident, ok := n.(*ast.Ident); ok && ident.Obj != nil {
				ident.Name = ""
				ident.Obj = nil
			}
			return true
		})
	}

	clearParameterNamesAndObjects(a)
	clearParameterNamesAndObjects(b)

	clearFunctionBodyBoundIdentifiers(a)
	clearFunctionBodyBoundIdentifiers(b)

	a = stripAllPos(a).(*ast.FuncDecl)
	b = stripAllPos(b).(*ast.FuncDecl)

	return reflect.DeepEqual(a, b)
}

// TODO: naming
func stripAllPos(n ast.Node) ast.Node {
	ast.Inspect(n, func(n ast.Node) bool {
		if n == nil {
			return false
		}

		v := reflect.ValueOf(n).Elem() // access the struct behind the interface and pointer
		for i := 0; i < v.NumField(); i++ {
			f := v.Field(i)
			// token.Pos(0) is a Pos value with underlying integer value 0, representing an undefined position.
			pos := token.Pos(0)
			if f.Type() == reflect.TypeOf(pos) {
				f.Set(reflect.ValueOf(pos))
			}
		}
		return true
	})
	return n
}

// TODO: name with Must or is Must with err?
func getBothFunctions(t *testing.T, asts map[string]*ast.File) (*ast.FuncDecl, *ast.FuncDecl) {
	dI, dJ := &ast.FuncDecl{}, &ast.FuncDecl{}

	for i, astI := range asts {
		for j, astJ := range asts {
			if i == j {
				continue
			}

			for _, declI := range astI.Decls {
				fnI, okI := declI.(*ast.FuncDecl)
				if !okI {
					continue
				}
				for _, declJ := range astJ.Decls {
					fnJ, okJ := declJ.(*ast.FuncDecl)
					if !okJ {
						continue
					}
					return fnI, fnJ
				}
			}
		}
	}

	t.Errorf("fail test")
	return dI, dJ
}

// TODO: rename a.go to i.go? same for the function. Naming tests
func TestFindDuplicateFunctions(t *testing.T) {
	testCases := []struct {
		name       string
		files      map[string]string
		areTheSame bool
	}{
		{
			name: "a",
			files: map[string]string{
				"a.go": `package main
						func A() {}`,
				"b.go": `package main
						func B()  {  }`,
			},
			areTheSame: true,
		},
		{
			name: "b",
			files: map[string]string{
				"a.go": `package main
						import "log"
						func A(sI string) {
							log.Println("s: ", sI)
						}`,
				"b.go": `package main
						import "log"
						func B(sJ string) {
							log.Println("s: ", sJ)
						}`,
			},
			areTheSame: true,
		},
		{
			name: "c",
			files: map[string]string{
				"a.go": `package main
						import "log"
						func A(sI string) {
							newSI := modifyString(sI)
							log.Println("s: ", newSI)
						}`,
				"b.go": `package main
						import "log"
						func B(sJ string) {
							newSJ := modifyString(sJ)
							log.Println("s: ", newSJ)
						}`,
			},
			areTheSame: true,
		},
		{
			name: "d",
			files: map[string]string{
				"a.go": `package main
						import "log"
						func A(sI string) {
							newSI := modifyStringI(sI)
							log.Println("s: ", newSI)
						}`,
				"b.go": `package main
						import "log"
						func B(sJ string) {
							newSJ := modifyStringJ(sJ)
							log.Println("s: ", newSJ)
						}`,
			},
			areTheSame: false,
		},
		{
			name: "e",
			files: map[string]string{
				"a.go": `package main
						import "log"
						func A() {
							s := "s"
							log.Println("s: ", s)
						}`,
				"b.go": `package main
						import "log"
						func B() {
							s := "s"
							log.Println("s: ", s)
						}`,
			},
			areTheSame: true,
		},
		{
			name: "f",
			files: map[string]string{
				"a.go": `package main
						import "log"
						func A() {
							sI := "sI"
							log.Println("s: ", sI)
						}`,
				"b.go": `package main
						import "log"
						func B() {
							sI := "sJ"
							log.Println("s: ", sJ)
						}`,
			},
			areTheSame: false,
		},
		{
			name: "g",
			files: map[string]string{
				"a.go": `package main
						import "log"
						func A() string {
							s := "s"
							log.Println("s: ", s)
							return s
						}`,
				"b.go": `package main
						import "log"
						func B() string {
							s := "s"
							log.Println("s: ", s)
							return s
						}`,
			},
			areTheSame: true,
		},
		{
			name: "h",
			files: map[string]string{
				"a.go": `package main
						import "log"
						func A() string {
							s := "s"
							log.Println("s: ", s)
							return s
						}`,
				"b.go": `package main
						import "log"
						func B() {
							s := "s"
							log.Println("s: ", s)
						}`,
			},
			areTheSame: false,
		},
		{
			name: "i",
			files: map[string]string{
				"a.go": `package main
						import "log"
						func A() {
							s := "s"
							log.Println("s: ", s)
						}`,
				"b.go": `package main
						import "log"
						func B() {
							s := 1
							log.Println("s: ", s)
						}`,
			},
			areTheSame: false,
		},
		{
			name: "j",
			files: map[string]string{
				"a.go": `package main
						import "log"
						func A() {
							s := "s"
							log.Println("s: ", s)
						}`,
				"b.go": `package main
						import "log"
						func B() {
							s := "s"
							s = ""
							s = "s"
							log.Println("s: ", s)
						}`,
			},
			areTheSame: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			functionI, functionJ := getBothFunctions(t, tMustCreateAsts(t, tc.files))
			if areFunctionsStructurallyEqual(functionI, functionJ) != tc.areTheSame {
				t.Errorf("fail test")
			}
		})
	}
}

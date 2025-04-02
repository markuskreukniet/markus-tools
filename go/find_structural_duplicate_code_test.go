package main

import (
	"go/ast"
	"go/token"
	"testing"

	"github.com/markuskreukniet/markus-tools/go/utils"
)

func tMustCreateAsts(t *testing.T, set *token.FileSet, files map[string]string) map[string]*ast.File {
	result, err := createAsts(set, files)
	return utils.TMust(t, result, err)
}

func tMustFindStructuralDuplicateFunctionParts(
	t *testing.T, set *token.FileSet, declI, declJ *ast.FuncDecl,
) [][]codeLocation {
	result, err := findStructuralDuplicateFunctionParts(set, declI, declJ)
	return utils.TMust(t, result, err)
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

	t.Errorf("fail test") // TODO: better fail?
	return dI, dJ
}

// TODO: Naming tests. Add tests with lower case function. Rename I an J postfixes?
func TestFindStructuralDuplicateFunctionBodyParts(t *testing.T) {
	testCases := []struct {
		name               string
		files              map[string]string
		numberOfDuplicates int // TODO: should become a duplicateCodeParts slice
	}{
		{
			name: "a",
			files: map[string]string{
				"a.go": `package main
						func A() string {
							s := "s"
							return s
						}`,
				"b.go": `package main
						func B() string  {
							s := "s"

							return s
						}`,
			},
			numberOfDuplicates: 1,
		},
		{
			name: "b",
			files: map[string]string{
				"a.go": `package main
						func A() string {
							return "sA"
						}`,
				"b.go": `package main
						func B() string   {
							return "sB"
						}`,
			},
			numberOfDuplicates: 0,
		},
		{
			name: "f",
			files: map[string]string{
				"a.go": `package main
								import "log"
								func A() {
									getSA := func() string {
										return "s"
									}
									log.Println("s: ", getSA())
								}`,
				"b.go": `package main
								import "log"
								func B() {
									getSB := func() string {
										return "s"
									}

									log.Println("s: ", getSB())
								}`,
			},
			numberOfDuplicates: 1,
		},
		// {
		// 	name: "g",
		// 	files: map[string]string{
		// 		"a.go": `package main
		// 				import "log"
		// 				func A() {
		// 					getSA := func() {
		// 						return "s"
		// 					}

		// 					log.Println("sA")

		// 					log.Println("s: ", getSA())
		// 				}`,
		// 		"b.go": `package main
		// 				import "log"
		// 				func B() {
		// 					getSB := func() {
		// 						return "s"
		// 					}

		// 					log.Println("sB")
		// 					log.Println("s: ", getSB())
		// 				}`,
		// 	},
		// 	numberOfDuplicates: 2,
		// },
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			set := token.NewFileSet()
			functionI, functionJ := getBothFunctions(t, tMustCreateAsts(t, set, tc.files))
			parts := tMustFindStructuralDuplicateFunctionParts(t, set, functionI, functionJ)

			if len(parts) != tc.numberOfDuplicates {
				t.Errorf("fail test") // TODO: better fail?
			}
		})
	}
}

// TODO: Naming tests. Add tests with lower case function. Rename I an J postfixes?
func TestFindStructuralDuplicateFunctions(t *testing.T) {
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
		// {
		// 	name: "d",
		// 	files: map[string]string{
		// 		"a.go": `package main
		// 				import "log"
		// 				func A(sI string) {
		// 					newSI := modifyStringI(sI)
		// 					log.Println("s: ", newSI)
		// 				}`,
		// 		"b.go": `package main
		// 				import "log"
		// 				func B(sJ string) {
		// 					newSJ := modifyStringJ(sJ)
		// 					log.Println("s: ", newSJ)
		// 				}`,
		// 	},
		// 	areTheSame: false,
		// },
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
		{
			name: "k",
			files: map[string]string{
				"a.go": `package main
						import "log"
						func A() {
							getSA := func() string {return "s"}
							log.Println("s: ", getSA())
						}`,
				"b.go": `package main
						import "log"
						func B() {
							getSB := func() string {
								return "s"
							}

							log.Println("s: ", getSB())
						}`,
			},
			areTheSame: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			functionI, functionJ := getBothFunctions(t, tMustCreateAsts(t, token.NewFileSet(), tc.files))
			if areFunctionsStructurallyEqual(functionI, functionJ) != tc.areTheSame {
				t.Errorf("fail test") // TODO: better fail?
			}
		})
	}
}

package main

import (
	"go/ast"
	"go/token"
	"strings"
	"testing"

	"github.com/markuskreukniet/markus-tools/go/utils"
)

func tMustCreateASTs(t *testing.T, set *token.FileSet, files map[string]string) map[string]*ast.File {
	result, err := createASTs(set, files)
	return utils.TMust(t, result, err)
}

func tMustFindDuplicateCodeRegions(t *testing.T, set *token.FileSet, declI, declJ *ast.FuncDecl) string {
	result, err := findDuplicateCodeRegions(set, declI, declJ)
	return utils.TMust(t, result, err)
}

func tMustFindFunctionallyEqualFunctions(t *testing.T, set *token.FileSet, asts map[string]*ast.File) string {
	result, err := findFunctionallyEqualFunctions(set, asts)
	return utils.TMust(t, result, err)
}

func tMustExtractNormalizedFunctions(
	t *testing.T, set *token.FileSet, asts map[string]*ast.File,
) []normalizedFunctionFile {
	result, err := extractNormalizedFunctions(set, asts)
	return utils.TMust(t, result, err)
}

func TestFindStructuralDuplicateFunctionBodyParts(t *testing.T) {
	testCases := []struct {
		name               string
		files              map[string]string
		numberOfDuplicates int // TODO: should become a string
	}{
		{
			name: "Different function names and positions",
			files: map[string]string{
				"i.go": `package main
						func I() string {
							s := "s"
							return s
						}`,
				"j.go": `package main
						func J() string  {
							s := "s"

							return s
						}`,
			},
			numberOfDuplicates: 1,
		},
		{
			name: "Different lowercase function names and positions",
			files: map[string]string{
				"i.go": `package main
						func i() string {
							s := "s"
							return s
						}`,
				"j.go": `package main
						func j() string  {
							s := "s"

							return s
						}`,
			},
			numberOfDuplicates: 1,
		},
		{
			name: "Different function names, positions, and return values",
			files: map[string]string{
				"i.go": `package main
						func I() string {
							return "sI"
						}`,
				"j.go": `package main
						func J() string   {
							return "sJ"
						}`,
			},
			numberOfDuplicates: 0,
		},
		{
			name: "Different function names, parameters order, positions, and return values order",
			files: map[string]string{
				"i.go": `package main
						func I(i int, s string) (int, string) {
							return "sI"
						}`,
				"j.go": `package main
						func J(s string, i int) (string, int)   {
							return "sJ"
						}`,
			},
			numberOfDuplicates: 0,
		},
		{
			name: "Different function names, positions, and the same 'Println' with a function call",
			files: map[string]string{
				"i.go": `package main
						import "log"
						func I() {
							getSI := func() string {
								return "s"
							}
							log.Println("s: ", getSI())
						}`,
				"j.go": `package main
						import "log"
						func J() {
							getSJ := func() string {
								return "s"
							}

							log.Println("s: ", getSJ())
						}`,
			},
			numberOfDuplicates: 1,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			set := token.NewFileSet()
			files := tMustExtractNormalizedFunctions(t, set, tMustCreateASTs(t, set, tc.files))

			// TODO: splitting on "\n\n" is duplicate code, but it is a temp solution
			parts := strings.Split(
				tMustFindDuplicateCodeRegions(t, set, files[0].functions[0].function, files[1].functions[0].function),
				"\n\n",
			)
			if (len(parts) != tc.numberOfDuplicates) && (parts[0] == "" && tc.numberOfDuplicates != 0) {
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
				"i.go": `package main
						func A() {}`,
				"j.go": `package main
						func B()  {  }`,
			},
			areTheSame: true,
		},
		{
			name: "b",
			files: map[string]string{
				"i.go": `package main
						import "log"
						func A(sI string) {
							log.Println("s: ", sI)
						}`,
				"j.go": `package main
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
				"i.go": `package main
						import "log"
						func A(sI string) {
							newSI := modifyString(sI)
							log.Println("s: ", newSI)
						}`,
				"j.go": `package main
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
				"i.go": `package main
						import "log"
						func A(sI string) {
							newSI := modifyStringI(sI)
							log.Println("s: ", newSI)
						}`,
				"j.go": `package main
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
				"i.go": `package main
						import "log"
						func A() {
							s := "s"
							log.Println("s: ", s)
						}`,
				"j.go": `package main
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
				"i.go": `package main
						import "log"
						func A() {
							sI := "sI"
							log.Println("s: ", sI)
						}`,
				"j.go": `package main
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
				"i.go": `package main
						import "log"
						func A() string {
							s := "s"
							log.Println("s: ", s)
							return s
						}`,
				"j.go": `package main
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
				"i.go": `package main
						import "log"
						func A() string {
							s := "s"
							log.Println("s: ", s)
							return s
						}`,
				"j.go": `package main
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
				"i.go": `package main
						import "log"
						func A() {
							s := "s"
							log.Println("s: ", s)
						}`,
				"j.go": `package main
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
				"i.go": `package main
						import "log"
						func A() {
							s := "s"
							log.Println("s: ", s)
						}`,
				"j.go": `package main
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
		// {
		// 	name: "k",
		// 	files: map[string]string{
		// 		"i.go": `package main
		// 				import "log"
		// 				func A() {
		// 					getSA := func() string {return "s"}
		// 					log.Println("s: ", getSA())
		// 				}`,
		// 		"j.go": `package main
		// 				import "log"
		// 				func B() {
		// 					getSB := func() string {
		// 						return "s"
		// 					}

		// 					log.Println("s: ", getSB())
		// 				}`,
		// 	},
		// 	areTheSame: true,
		// },
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			set := token.NewFileSet()
			utils.TMustAssertEqualBools(
				t,
				!utils.IsBlank(tMustFindFunctionallyEqualFunctions(t, set, tMustCreateASTs(t, set, tc.files))),
				tc.areTheSame,
			)
		})
	}
}

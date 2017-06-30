package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
)

func main() {
	pDir := flag.String("dir", "", "the directory")
	pAst := flag.Bool("ast", false, "print ast")
	flag.Parse()

	if *pDir == "" {
		fmt.Printf("Error: -dir flag is required\n")
		flag.Usage()
	}

	files, err := filepath.Glob(fmt.Sprintf("%s/*.go", *pDir))
	if err != nil {
		fmt.Printf("Could not read directory %s : %v\n", *pDir, err)
		os.Exit(-1)
	}

	// token store
	fset := token.NewFileSet()

	types := make([]string, 0)
	for _, file := range files {
		// include comments in the syntax tree
		f, err := parser.ParseFile(fset, file, nil, parser.ParseComments)
		if err != nil {
			fmt.Printf("Error parsing file %s : %v\n", file, err)
			os.Exit(-1)
		}

		// Print the first file's AST and exit
		// useful for debugging purpose
		if *pAst {
			ast.Print(fset, f)
			os.Exit(0)
		}

		// Parse declarations
		for _, decl := range f.Decls {
			if gDecl, gDeclOk := decl.(*ast.GenDecl); gDeclOk {
				for _, spec := range gDecl.Specs {
					// If it is type spec
					if tSpec, tSpecOk := spec.(*ast.TypeSpec); tSpecOk {
						// Are we looking at a struct type declrataion?
						// then we want it
						if _, stTypeOk := tSpec.Type.(*ast.StructType); stTypeOk {
							types = append(types, tSpec.Name.Name)
						}
					}
				}
			}
		}
	}

	// convert to json
	oBytes, err := json.MarshalIndent(types, "", "    ")
	if err != nil {
		fmt.Printf("Error converting to json : %v\n", err)
		os.Exit(-1)
	}
	fmt.Println(string(oBytes))
}

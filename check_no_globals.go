package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
)

func checkNoGlobals(path string) ([]string, error) {
	messages := []string{}

	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			return nil
		}

		fset := token.NewFileSet()
		pkgs, err := parser.ParseDir(fset, path, nil, 0)
		if err != nil {
			return err
		}

		for _, pkg := range pkgs {
			for _, file := range pkg.Files {
				for _, decl := range file.Decls {
					genDecl, ok := decl.(*ast.GenDecl)
					if !ok {
						continue
					}
					if genDecl.Tok != token.VAR {
						continue
					}
					filename := fset.Position(genDecl.TokPos).Filename
					line := fset.Position(genDecl.TokPos).Line
					valueSpec := genDecl.Specs[0].(*ast.ValueSpec)
					for i := 0; i < len(valueSpec.Names); i++ {
						name := valueSpec.Names[i].Name
						if name == "_" {
							continue
						}
						message := fmt.Sprintf("%s:%d %s is a global variable", filename, line, name)
						messages = append(messages, message)
					}
				}
			}
		}
		return nil
	})

	return messages, err
}
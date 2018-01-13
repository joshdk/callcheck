// Copyright 2018 Josh Komoroske. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE.txt file.

package graph

import (
	"go/ast"

	"golang.org/x/tools/go/loader"
)

type FuncDecl struct {
	Name     string
	Package  string
	Position string
	Calls    []FuncCall
}

type FuncCall struct {
	Name     string
	Package  string
	Position string
}

func Program(program *loader.Program) (map[string]FuncDecl, error) {
	decls := make(map[string]FuncDecl)

	// Check every package that belongs to the program.
	for _, pkgInfo := range program.AllPackages {
		// Check every file that belongs to the package.
		for _, astFile := range pkgInfo.Files {
			// Check every function declared in the file.
			for _, f := range astFile.Decls {
				fn, ok := f.(*ast.FuncDecl)
				if !ok {
					continue
				}

				// Attempt to fully qualify the function declaration name and package.
				pkgName, name, ok := Qualify(pkgInfo, fn)
				if !ok {
					continue
				}

				decls[name] = FuncDecl{
					Name:     name,
					Package:  pkgName,
					Position: program.Fset.Position(fn.Pos()).String(),
					Calls:    []FuncCall{},
				}

				vis := funcDeclVisitor{
					pkgInfo,
					program.Fset,
					name,
					decls,
				}

				// Walk contents of the function declaration
				ast.Walk(&vis, fn)
			}
		}
	}

	return decls, nil
}

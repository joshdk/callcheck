// Copyright 2018 Josh Komoroske. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE.txt file.

package graph

import (
	"go/ast"
	"go/token"

	"golang.org/x/tools/go/loader"
)

type funcDeclVisitor struct {
	pkg     *loader.PackageInfo
	fset    *token.FileSet
	current string
	decls   map[string]FuncDecl
}

// Visit is intended to traverses the contents of an ast.FuncDecl, and will
// record the existence of all function calls located within the function body.
func (v *funcDeclVisitor) Visit(node ast.Node) ast.Visitor {

	// The visitor is only concerned with function calls. If the current node
	// is not a CallExpr, then no additional processing is done.
	stmt, ok := node.(*ast.CallExpr)
	if !ok {
		return v
	}

	// Attempt to fully qualify the function call name and package.
	if pkgName, funcName, ok := Qualify(v.pkg, stmt); ok {

		call := FuncCall{
			Name:     funcName,
			Package:  pkgName,
			Position: v.fset.Position(stmt.Pos()).String(),
		}

		// Record that this function call exists inside the parent function
		// body.
		v.add(call)
	}

	return v
}

func (v *funcDeclVisitor) add(call FuncCall) {
	curr := v.decls[v.current]
	curr.Calls = append(curr.Calls, call)
	v.decls[v.current] = curr
}

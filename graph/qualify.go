// Copyright 2018 Josh Komoroske. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE.txt file.

package graph

import (
	"go/ast"
	"go/types"

	"golang.org/x/tools/go/loader"
)

func Qualify(pkg *loader.PackageInfo, node ast.Node) (string, string, bool) {
	var ident *ast.Ident

	switch kind := node.(type) {
	// Match function declarations
	case *ast.FuncDecl:
		ident = kind.Name

		// Match function invocations
	case *ast.CallExpr:
		switch expr := kind.Fun.(type) {
		case *ast.SelectorExpr:
			ident = expr.Sel
		case *ast.Ident:
			ident = expr
		}
	}

	switch fn := pkg.ObjectOf(ident).(type) {
	case *types.Func:
		// Builtin interface function call like err.Error()
		if fn.Pkg() == nil {
			return "", fn.FullName(), true
		}
		return fn.Pkg().Path(), fn.FullName(), true
	case *types.Builtin:
		return "", fn.Name(), true
	default:
		return "", "", false
	}
}

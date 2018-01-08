// Copyright 2018 Josh Komoroske. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE.txt file.

package cmd

import (
	"fmt"
	"go/build"

	"github.com/kisielk/gotool"
	"golang.org/x/tools/go/loader"

	"github.com/joshdk/callcheck/graph"
)

func Cmd(args []string) error {
	paths := gotool.ImportPaths(args)

	cfg := loader.Config{
		Build: &build.Default,
	}

	_, err := cfg.FromArgs(paths, false)
	if err != nil {
		return err
	}

	program, err := cfg.Load()
	if err != nil {
		return err
	}

	decls, err := graph.Program(program)
	if err != nil {
		return err
	}

	for decl, calls := range decls {
		fmt.Printf("%s | %s\n", decl.Position, decl.Name)
		for _, call := range calls {
			fmt.Printf(" └─ %s | %s\n", call.Position, call.Name)
		}
	}

	return nil
}

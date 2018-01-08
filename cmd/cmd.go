// Copyright 2018 Josh Komoroske. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE.txt file.

package cmd

import (
	"fmt"
	"go/build"

	"github.com/kisielk/gotool"
	"golang.org/x/tools/go/loader"
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

	for _, pkgInfo := range program.AllPackages {
		fmt.Println(pkgInfo.Pkg.Path())
	}

	return nil
}

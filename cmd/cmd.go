// Copyright 2018 Josh Komoroske. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE.txt file.

package cmd

import (
	"errors"
	"go/build"

	"github.com/kisielk/gotool"
	"golang.org/x/tools/go/loader"

	"github.com/joshdk/callcheck/config"
	"github.com/joshdk/callcheck/graph"
)

func Cmd(args []string) error {
	if len(args) == 0 {
		args = []string{"./..."}
	}

	checkCfg, err := config.Load()
	if err != nil {
		return err
	}

	paths := gotool.ImportPaths(args)

	cfg := loader.Config{
		Build: &build.Default,
	}

	if _, err := cfg.FromArgs(paths, false); err != nil {
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

	if violations(decls, checkCfg) {
		return errors.New("policy violations found")
	}

	return nil
}

// Copyright 2018 Josh Komoroske. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE.txt file.

package main

import (
	"os"

	"github.com/joshdk/callcheck/cmd"
)

func main() {
	err := cmd.Cmd(os.Args[1:])
	if err != nil {
		os.Exit(1)
	}
}

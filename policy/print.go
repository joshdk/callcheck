// Copyright 2018 Josh Komoroske. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE.txt file.

package policy

import (
	"fmt"
	"strings"
)

func (decl Decl) String() string {
	var (
		results = decl.string("")
		buffer  string
		max     int
	)

	// Find the width of the longest position text. This will be used to
	// calculate padding, so that the columns can be aligned vertically.
	for _, result := range results {
		if len(result[0]) > max {
			max = len(result[0])
		}
	}

	// Format each line, with correct column padding.
	for _, result := range results {
		buffer += fmt.Sprintf("%s%s → %s\n", result[0], strings.Repeat(" ", max-len(result[0])), result[1])
	}

	return buffer
}

func (decl Decl) string(prefix string) [][2]string {
	results := [][2]string{{
		fmtUnknown(decl.Position),
		prefix + fmtUnknown(decl.Name),
	}}

	for index, call := range decl.Calls {
		suffix := "├── "
		runner := "│   "
		if index == len(decl.Calls)-1 {
			suffix = "└── "
			runner = "    "
		}

		results = append(results, [2]string{
			fmtUnknown(call.Position),
			prefix + suffix + fmtUnknown(call.Name),
		})

		results = append(results, call.Decl.string(prefix+runner)...)
	}

	return results
}

func fmtUnknown(name string) string {
	if name == "" {
		return "???"
	}
	return name
}

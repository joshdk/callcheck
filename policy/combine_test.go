// Copyright 2018 Josh Komoroske. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE.txt file.

package policy

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCombine(t *testing.T) {

	tests := []struct {
		title    string
		first    Decl
		second   Decl
		match    string
		split    string
		expected string
		err      string
	}{
		{
			title:    "identical single decls",
			first:    Decl{Name: "main", Position: "main.go:0"},
			second:   Decl{Name: "main", Position: "main.go:0"},
			match:    "main",
			split:    "",
			expected: "main.go:0 → main",
		},
		{
			title:  "different single decls",
			first:  Decl{Name: "main", Position: "main.go:0"},
			second: Decl{Name: "init", Position: "init.go:0"},
			match:  "main",
			err:    "declarations required to match but did not",
		},
		{
			title: "decls that do not split",
			first: Decl{
				Name:     "main",
				Position: "main.go:0",
				Calls: []Call{
					{
						Name:     "run",
						Position: "main.go:1",
						Index:    0,
						Decl: Decl{
							Name:     "run",
							Position: "run.go:0",
						},
					},
				},
			},
			second: Decl{
				Name:     "main",
				Position: "main.go:0",
				Calls: []Call{
					{
						Name:     "run",
						Position: "main.go:1",
						Index:    0,
						Decl: Decl{
							Name:     "run",
							Position: "run.go:0",
						},
					},
				},
			},
			match: "main",
			split: "run",
			err:   "declarations required to not match but did",
		},
		{
			title: "reversed order decls that do not split",
			first: Decl{
				Name:     "main",
				Position: "main.go:0",
				Calls: []Call{
					{
						Name:     "run",
						Position: "main.go:1",
						Index:    0,
						Decl: Decl{
							Name:     "run",
							Position: "run.go:0",
							Calls: []Call{
								{
									Name:     "load",
									Position: "run.go:1",
									Index:    0,
									Decl: Decl{
										Name:     "load",
										Position: "load.go:0",
									},
								},
							},
						},
					},
				},
			},
			second: Decl{
				Name:     "main",
				Position: "main.go:0",
				Calls: []Call{
					{
						Name:     "run",
						Position: "main.go:1",
						Index:    0,
						Decl: Decl{
							Name:     "run",
							Position: "run.go:0",
						},
					},
				},
			},
			match: "main",
			split: "run",
			err:   "declarations required to not match but did",
		},
	}

	for index, test := range tests {
		name := fmt.Sprintf("#%d - %s", index, test.title)

		t.Run(name, func(t *testing.T) {
			actual, err := combineDecls(test.first, test.second, test.match, test.split)

			if test.err != "" {
				require.EqualError(t, err, test.err)
			} else {
				assert.Equal(t, strings.TrimSpace(unindent(test.expected)), strings.TrimSpace(actual.String()))
			}
		})
	}
}

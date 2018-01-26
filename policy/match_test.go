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

	"github.com/joshdk/callcheck/graph"
)

func TestWrapDecl(t *testing.T) {

	tests := []struct {
		title    string
		wrapper  Decl
		wrapped  Decl
		expected string
	}{
		{
			title:    "same node",
			wrapper:  Decl{Name: "main", Position: "main.go:0"},
			wrapped:  Decl{Name: "main", Position: "main.go:0"},
			expected: `main.go:0 → main`,
		},
		{
			title: "single-node wrapper with multi-node wrapped",
			wrapper: Decl{
				Name:     "main",
				Position: "main.go:0",
			},
			wrapped: Decl{
				Name:     "main",
				Position: "main.go:0",
				Calls: []Call{
					{
						Name:     "load",
						Position: "main.go:1",
						Index:    1,
						Decl: Decl{
							Name:     "load",
							Position: "load.go:0",
						},
					},
				},
			},
			expected: `
				main.go:0 → main
				main.go:1 → └── load
				load.go:0 →     load
			`,
		},
		{
			title:   "single-node wrapper with multi-node branched wrapped",
			wrapper: Decl{Name: "main"},
			wrapped: Decl{
				Name:     "main",
				Position: "main.go:0",
				Calls: []Call{
					{
						Name:     "load",
						Position: "main.go:1",
						Index:    1,
						Decl: Decl{
							Name:     "load",
							Position: "load.go:0",
						},
					},
					{
						Name:     "run",
						Index:    2,
						Position: "main.go:2",
						Decl: Decl{
							Name:     "run",
							Position: "run.go:0",
						},
					},
				},
			},
			expected: `
				main.go:0 → main
				main.go:1 → ├── load
				load.go:0 → │   load
				main.go:2 → └── run
				run.go:0  →     run
			`,
		},
		{
			title: "multi-node wrapper with single-node wrapped",
			wrapper: Decl{
				Name:     "main",
				Position: "main.go:0",
				Calls: []Call{
					{
						Name:     "load",
						Position: "main.go:1",
						Index:    1,
						Decl: Decl{
							Name: "load",
						},
					},
				},
			},
			wrapped: Decl{
				Name:     "load",
				Position: "load.go:0",
			},
			expected: `
				main.go:0 → main
				main.go:1 → └── load
				load.go:0 →     load
			`,
		},
		{
			title: "multi-node wrapper with multi-node wrapped",
			wrapper: Decl{
				Name:     "main",
				Position: "main.go:0",
				Calls: []Call{
					{
						Name:     "load",
						Position: "main.go:1",
						Index:    1,
						Decl: Decl{
							Name:     "load",
							Position: "load.go:0",
						},
					},
				},
			},
			wrapped: Decl{
				Name:     "load",
				Position: "load.go:0",
				Calls: []Call{
					{
						Name:     "read",
						Position: "load.go:1",
						Index:    1,
						Decl: Decl{
							Name:     "read",
							Position: "read.go:0",
						},
					},
				},
			},
			expected: `
				main.go:0 → main
				main.go:1 → └── load
				load.go:0 →     load
				load.go:1 →     └── read
				read.go:0 →         read
			`,
		},
		{
			title: "multi-node wrapper with multi-node branched wrapped",
			wrapper: Decl{
				Name:     "main",
				Position: "main.go:0",
				Calls: []Call{
					{
						Name:     "load",
						Position: "main.go:1",
						Index:    1,
						Decl: Decl{
							Name:     "load",
							Position: "load.go:0",
						},
					},
				},
			},
			wrapped: Decl{
				Name:     "load",
				Position: "load.go:0",
				Calls: []Call{
					{
						Name:     "read",
						Position: "load.go:1",
						Index:    1,
						Decl: Decl{
							Name:     "read",
							Position: "read.go:0",
						},
					},
					{
						Name:     "parse",
						Position: "load.go:2",
						Index:    2,
						Decl: Decl{
							Name:     "parse",
							Position: "parse.go:0",
						},
					},
				},
			},
			expected: `
				main.go:0  → main
				main.go:1  → └── load
				load.go:0  →     load
				load.go:1  →     ├── read
				read.go:0  →     │   read
				load.go:2  →     └── parse
				parse.go:0 →         parse
			`,
		},
	}

	for index, test := range tests {
		name := fmt.Sprintf("#%d - %s", index, test.title)

		t.Run(name, func(t *testing.T) {
			actual := wrapDecl(test.wrapper, test.wrapped)
			assert.Equal(t, strings.TrimSpace(unindent(test.expected)), strings.TrimSpace(actual.String()))
		})
	}
}

func TestWalk(t *testing.T) {

	cg := map[string]graph.FuncDecl{
		"main": {
			Name:     "main",
			Position: "main.go:0",
			Calls: []graph.FuncCall{
				{Name: "parse", Position: "main.go:1"},
				{Name: "run", Position: "main.go:2"},
			},
		},
		"parse": {
			Name:     "parse",
			Position: "parse.go:0",
			Calls: []graph.FuncCall{
				{Name: "load", Position: "parse.go:1"},
				{Name: "panic", Position: "parse.go:2"},
			},
		},
		"run": {
			Name:     "run",
			Position: "run.go:0",
			Calls: []graph.FuncCall{
				{Name: "exec", Position: "run.go:1"},
				{Name: "panic", Position: "run.go:2"},
			},
		},
		"load": {
			Name:     "load",
			Position: "load.go:0",
			Calls: []graph.FuncCall{
				{Name: "panic", Position: "load.go:1"},
			},
		},
		"exec": {
			Name:     "exec",
			Position: "exec.go:0",
			Calls: []graph.FuncCall{
				{Name: "panic", Position: "exec.go:1"},
			},
		},
		"panic": {
			Name:     "panic",
			Position: "panic.go:0",
		},
		"recurse": {
			Name:     "recurse",
			Position: "recurse.go:0",
			Calls: []graph.FuncCall{
				{Name: "recurse", Position: "recurse.go:1"},
				{Name: "corecurse", Position: "recurse.go:2"},
			},
		},
		"corecurse": {
			Name:     "corecurse",
			Position: "corecurse.go:0",
			Calls: []graph.FuncCall{
				{Name: "recurse", Position: "corecurse.go:1"},
			},
		},
	}

	tests := []struct {
		title string
		start string
		end   string
		graph map[string]graph.FuncDecl
		paths []string
	}{
		{
			title: "nil graph",
		},
		{
			title: "node to same node",
			start: "main",
			end:   "main",
			graph: cg,
			paths: []string{"main.go:0 → main"},
		},
		{
			title: "recursive",
			start: "recurse",
			end:   "recurse",
			graph: cg,
			paths: []string{"recurse.go:0 → recurse"},
		},
		{
			title: "corecursive",
			start: "recurse",
			end:   "corecurse",
			graph: cg,
			paths: []string{
				`
					recurse.go:0   → recurse
					recurse.go:2   → └── corecurse
					corecurse.go:0 →     corecurse
				`,
			},
		},
		{
			title: "main > exec",
			start: "main",
			end:   "exec",
			graph: cg,
			paths: []string{
				`
					main.go:0 → main
					main.go:2 → └── run
					run.go:0  →     run
					run.go:1  →     └── exec
					exec.go:0 →         exec
				`,
			},
		},
		{
			title: "main > panic",
			start: "main",
			end:   "panic",
			graph: cg,
			paths: []string{
				`
					main.go:0  → main
					main.go:1  → └── parse
					parse.go:0 →     parse
					parse.go:1 →     └── load
					load.go:0  →         load
					load.go:1  →         └── panic
					panic.go:0 →             panic
				`,
				`
					main.go:0  → main
					main.go:1  → └── parse
					parse.go:0 →     parse
					parse.go:2 →     └── panic
					panic.go:0 →         panic

				`,
				`
					main.go:0  → main
					main.go:2  → └── run
					run.go:0   →     run
					run.go:1   →     └── exec
					exec.go:0  →         exec
					exec.go:1  →         └── panic
					panic.go:0 →             panic
				`,
				`
					main.go:0  → main
					main.go:2  → └── run
					run.go:0   →     run
					run.go:2   →     └── panic
					panic.go:0 →         panic
				`,
			},
		},
	}

	for index, test := range tests {
		name := fmt.Sprintf("#%d - %s", index, test.title)

		t.Run(name, func(t *testing.T) {
			actual := walker(test.start, test.end, test.graph)
			require.Equal(t, len(test.paths), len(actual))

			for index, path := range actual {
				assert.Equal(t, strings.TrimSpace(unindent(test.paths[index])), strings.TrimSpace(path.String()))
			}
		})
	}
}

func unindent(body string) string {
	return strings.Replace(body, "\t", "", -1)
}

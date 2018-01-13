// Copyright 2018 Josh Komoroske. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE.txt file.

package policy

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/joshdk/callcheck/graph"
)

func TestMatch(t *testing.T) {

	type test struct {
		name    string
		graph   map[string]graph.FuncDecl
		root    string
		matches bool
	}

	suites := []struct {
		name   string
		policy Policy
		tests  []test
	}{

		// Policy that is empty (has a nil rule)
		{
			name: "empty",
			policy: Policy{
				Name:        "empty",
				Description: "Nothing to see here, move along...",
			},
			tests: []test{
				{
					name:    "nil graph",
					matches: false,
					graph:   nil,
				},
				{
					name:    "only panic",
					matches: false,
					graph: map[string]graph.FuncDecl{
						"panic": {Name: "panic"},
					},
					root: "panic",
				},
			},
		},

		// Policy that matches calls to panic()
		{
			name: "forbid-panic",
			policy: Policy{
				Name:        "forbid-panic",
				Description: "DON'T PANIC",
				Rule: &Node{
					Name: "panic",
				},
			},
			tests: []test{
				{
					name:    "nil graph",
					matches: false,
					graph:   nil,
				},
				{
					name:    "recursive",
					matches: false,
					graph: map[string]graph.FuncDecl{
						"recurse": {
							Name: "recurse",
							Calls: []graph.FuncCall{
								{Name: "recurse"},
							},
						},
					},
					root: "recurse",
				},
				{
					name:    "only panic",
					matches: true,
					graph: map[string]graph.FuncDecl{
						"panic": {Name: "panic"},
					},
					root: "panic",
				},
				{
					name:    "only recover",
					matches: false,
					graph: map[string]graph.FuncDecl{
						"recover": {Name: "recover"},
					},
					root: "recover",
				},
				{
					name:    "main > panic",
					matches: true,
					graph: map[string]graph.FuncDecl{
						"panic": {Name: "panic"},
						"main": {
							Name: "main",
							Calls: []graph.FuncCall{
								{Name: "panic"},
							},
						},
					},
					root: "main",
				},
				{
					name:    "main > run > panic",
					matches: true,
					graph: map[string]graph.FuncDecl{
						"panic": {Name: "panic"},
						"main": {
							Name: "main",
							Calls: []graph.FuncCall{
								{Name: "parse"},
								{Name: "run"},
								{Name: "log"},
							},
						},
						"run": {
							Name: "run",
							Calls: []graph.FuncCall{
								{Name: "check"},
								{Name: "panic"},
							},
						},
					},
					root: "main",
				},
			},
		},

		// Policy that matches a recursive calls()
		{
			name: "forbid-recurse-recurse",
			policy: Policy{
				Name:        "forbid-recurse-recurse",
				Description: "Forbid recurse calling itself",
				Rule: &Node{
					Name: "recurse",
					Calls: []*Node{
						{Name: "recurse"},
					},
				},
			},
			tests: []test{
				{
					name:    "nil graph",
					matches: false,
					graph:   nil,
				},
				{
					name:    "recursive",
					matches: true,
					graph: map[string]graph.FuncDecl{
						"recurse": {
							Name: "recurse",
							Calls: []graph.FuncCall{
								{Name: "recurse"},
							},
						},
					},
					root: "recurse",
				},
				{
					name:    "corecursive",
					matches: true,
					graph: map[string]graph.FuncDecl{
						"recurse": {
							Name: "recurse",
							Calls: []graph.FuncCall{
								{Name: "corecurse"},
							},
						},
						"corecurse": {
							Name: "corecurse",
							Calls: []graph.FuncCall{
								{Name: "recurse"},
							},
						},
					},
					root: "corecurse",
				},
			},
		},

		// Policy that matches calls to parse() followed by panic()
		{
			name: "forbid-parse-panic",
			policy: Policy{
				Name:        "forbid-parse-panic",
				Description: "Forbid parse calling panic",
				Rule: &Node{
					Name: "parse",
					Calls: []*Node{
						{Name: "panic"},
					},
				},
			},
			tests: []test{
				{
					name:    "nil graph",
					matches: false,
					graph:   nil,
				},
				{
					name:    "recursive",
					matches: false,
					graph: map[string]graph.FuncDecl{
						"recurse": {
							Name: "recurse",
							Calls: []graph.FuncCall{
								{Name: "recurse"},
							},
						},
					},
					root: "recurse",
				},
				{
					name:    "only panic",
					matches: false,
					graph: map[string]graph.FuncDecl{
						"panic": {},
					},
					root: "panic",
				},
				{
					name:    "parse > panic",
					matches: true,
					graph: map[string]graph.FuncDecl{
						"parse": {
							Name: "parse",
							Calls: []graph.FuncCall{
								{Name: "panic"},
							},
						},
						"panic": {},
					},
					root: "parse",
				},
				{
					name:    "main > parse > panic",
					matches: true,
					graph: map[string]graph.FuncDecl{
						"main": {
							Name: "main",
							Calls: []graph.FuncCall{
								{Name: "parse"},
							},
						},
						"parse": {
							Name: "parse",
							Calls: []graph.FuncCall{
								{Name: "panic"},
							},
						},
						"panic": {},
					},
					root: "main",
				},
				{
					name:    "main > parse > load > panic",
					matches: true,
					graph: map[string]graph.FuncDecl{
						"main": {
							Name: "main",
							Calls: []graph.FuncCall{
								{Name: "parse"},
							},
						},
						"parse": {
							Name: "parse",
							Calls: []graph.FuncCall{
								{Name: "load"},
							},
						},
						"load": {
							Name: "load",
							Calls: []graph.FuncCall{
								{Name: "panic"},
							},
						},
						"panic": {},
					},
					root: "main",
				},
			},
		},

		// Policy that branches out
		{
			name: "forbid-complex",
			policy: Policy{
				Name:        "forbid-complex",
				Description: "Forbid a complex structure",
				Rule: &Node{
					Name: "root",
					Calls: []*Node{
						{
							Name: "mid1",
							Calls: []*Node{
								{
									Name: "mid2",
									Calls: []*Node{
										{Name: "leaf1"},
									},
								},
							},
						},
						{
							Name: "mid3",
							Calls: []*Node{
								{
									Name: "mid4",
									Calls: []*Node{
										{Name: "leaf2"},
									},
								},
							},
						},
						{Name: "leaf3"},
					},
				},
			},
			tests: []test{
				{
					name:    "nil graph",
					matches: false,
					graph:   nil,
				},
				{
					name:    "recursive",
					matches: false,
					graph: map[string]graph.FuncDecl{
						"recurse": {
							Name: "recurse",
							Calls: []graph.FuncCall{
								{Name: "recurse"},
							},
						},
					},
					root: "recurse",
				},
				{
					name:    "exact",
					matches: true,
					graph: map[string]graph.FuncDecl{
						"root": {
							Name: "root",
							Calls: []graph.FuncCall{
								{Name: "mid1"},
								{Name: "mid3"},
								{Name: "leaf3"},
							},
						},
						"mid1": {
							Name: "mid1",
							Calls: []graph.FuncCall{
								{Name: "mid2"},
							},
						},
						"mid2": {
							Name: "mid2",
							Calls: []graph.FuncCall{
								{Name: "leaf1"},
							},
						},
						"mid3": {
							Name: "mid3",
							Calls: []graph.FuncCall{
								{Name: "mid4"},
							},
						},
						"mid4": {
							Name: "mid4",
							Calls: []graph.FuncCall{
								{Name: "leaf2"},
							},
						},
						"leaf3": {
							Name: "leaf3",
						},
					},
					root: "root",
				},
			},
		},
	}

	for suiteIndex, suite := range suites {
		for testIndex, test := range suite.tests {
			name := fmt.Sprintf("%s #%d > %s #%d", suite.name, suiteIndex, test.name, testIndex)

			t.Run(name, func(t *testing.T) {
				matched := IsForbidden(suite.policy, test.root, test.graph)
				assert.Equal(t, test.matches, matched)
			})
		}
	}

}

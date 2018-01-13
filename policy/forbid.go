// Copyright 2018 Josh Komoroske. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE.txt file.

package policy

import (
	"github.com/joshdk/callcheck/graph"
)

// IsForbidden checks if a policy fully matches a graph, rooted on a specific
// function. Will return true if all policy rule nodes are matched within the
// graph. Rules that are nil will always return false.
func IsForbidden(policy Policy, root string, graph map[string]graph.FuncDecl) bool {
	// If the node is nil, then nothing can be matched.
	if policy.Rule == nil {
		return false
	}

	// If the graph is nil, then nothing can be matched.
	if graph == nil {
		return false
	}

	initRule(policy.Rule)

	return checkForbidden(policy.Rule, root, graph)
}

func checkForbidden(node *Node, current string, graph map[string]graph.FuncDecl) bool {
	// If this specific node has already visited this specific function, then
	// do not try a second time.
	if _, found := node.Visited[current]; found {
		return false
	}

	// Record that this specific node has already visited this specific
	// function.
	node.Visited[current] = struct{}{}

	calls := graph[current].Calls

	// If this nodes name matches the current function...
	if node.Name == current {

		for _, ruleCall := range node.Calls {
			for {
				if len(calls) == 0 {
					return false
				}

				curr := calls[0]
				calls = calls[1:]

				if checkForbidden(ruleCall, curr.Name, graph) {
					break
				}
			}
		}

		return true
	}

	// If this nodes name does not match the current function...
	for _, call := range calls {
		if checkForbidden(node, call.Name, graph) {
			return true
		}
	}

	return false
}

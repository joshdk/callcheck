// Copyright 2018 Josh Komoroske. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE.txt file.

package policy

import (
	"errors"

	"github.com/joshdk/callcheck/graph"
)

type Decl struct {
	Position string
	Name     string
	Calls    []Call
}

type Call struct {
	Position string
	Name     string
	Decl     Decl
	Index    int
}

type Goal struct {
	startDecl string
	endCall   string
}

func MatchingPaths(graph map[string]graph.FuncDecl, policy Policy) []Decl {
	if policy.Rule == nil {
		return nil
	}

	if graph == nil {
		return nil
	}

	// Our rule does not have had any calls, and must be a decl all on its own.
	if len(policy.Rule.Calls) == 0 {
		// Check that our rule actually exist in the graph.
		if decl, found := graph[policy.Rule.Name]; found {
			return []Decl{{Name: decl.Name, Position: decl.Position}}
		}
		return nil
	}

	nodeToGoalMapping := splitRule(policy.Rule)
	goalToDeclsMapping := make(map[*Goal][]Decl)

	for _, goals := range nodeToGoalMapping {
		for _, goal := range goals {
			decls := walker(goal.startDecl, goal.endCall, graph)
			goalToDeclsMapping[goal] = decls
		}
	}

	return genMatches(policy.Rule, nodeToGoalMapping, goalToDeclsMapping)
}

// walker traverses the given call graph from the function named start and
// returns all distinct paths to the function named end. A value of nil is
// returned if no paths are found. All returned paths are guaranteed to be
// linear (do not branch).
func walker(start string, end string, graph map[string]graph.FuncDecl) []Decl {
	visited := make(map[string]struct{})
	return paths(start, end, visited, graph)
}

// paths is an internal function behind walker.
func paths(current string, end string, visited map[string]struct{}, graph map[string]graph.FuncDecl) []Decl {
	if graph == nil {
		return nil
	}

	me := Decl{
		Position: graph[current].Position,
		Name:     current,
	}

	if current == end {
		return []Decl{me}
	}

	if _, found := visited[current]; found {
		return nil
	}

	startDecl := graph[current]

	visited[current] = struct{}{}

	var results []Decl

	for index, call := range startDecl.Calls {
		paths := paths(call.Name, end, visited, graph)
		for _, path := range paths {
			results = append(results, Decl{
				Name:     current,
				Position: startDecl.Position,
				Calls: []Call{
					{
						call.Position,
						path.Name,
						path,
						index,
					},
				},
			})
		}
	}

	return results
}

func walkRule(node *Node, goals map[*Node][]*Goal) {
	if node == nil {
		return
	}

	for _, call := range node.Calls {
		if _, found := goals[node]; !found {
			goals[node] = []*Goal{}
		}

		goals[node] = append(goals[node], &Goal{node.Name, call.Name})
		walkRule(call, goals)
	}
}

// splitRule splits the given node into a list of goals for every sub-node.
func splitRule(node *Node) map[*Node][]*Goal {
	goals := make(map[*Node][]*Goal)
	walkRule(node, goals)
	return goals
}

func combineDecls(first Decl, second Decl, mustMatch string, mustSplit string) (Decl, error) {
	// Sanity check declarations.
	switch {
	case first.Name == "" || second.Name == "":
		panic("name is blank")

		// Nodes are supposed to match, but didn't
	case first.Name == mustMatch && second.Name != first.Name:
		return Decl{}, errors.New("declarations required to match but did not")

		// Nodes are not supposed to match, but did
	case first.Name == mustSplit && second.Name == first.Name:
		return Decl{}, errors.New("declarations required to not match but did")

	case len(second.Calls) > 1:
		panic("more than 1 call")

	case first.Name != second.Name:
		panic("disjoint declarations")

	}

	// Nodes are the same, merge
	merged, err := combineCalls(first.Calls, second.Calls, mustMatch, mustSplit)
	if err != nil {
		return Decl{}, err
	}

	return Decl{
		Name:     first.Name,
		Position: first.Position,
		Calls:    merged,
	}, nil
}

func combineCalls(first []Call, second []Call, mustMatch string, mustSplit string) ([]Call, error) {
	// Sanity check calls.
	switch {
	case len(first) == 0 && len(second) == 0:
		return nil, nil

	case len(second) >= 2:
		panic("more than 1 call")

	case len(first) == 0 && len(second) >= 1:
		panic("first had no calls but second had calls")

	case len(first) >= 1 && len(second) == 0:
		panic("first had calls but second had no calls")
	}

	firstCall, secondCall := first[0], second[0]

	// These two calls are the same, merge them.
	if firstCall.Name == secondCall.Name {
		merged, err := combineDecls(firstCall.Decl, secondCall.Decl, mustMatch, mustSplit)
		if err != nil {
			return nil, err
		}

		return []Call{{
			Name:     firstCall.Name,
			Position: firstCall.Position,
			Index:    firstCall.Index,
			Decl:     merged,
		}}, nil
	}

	// These two calls are not the same, check if they are ordered.
	if firstCall.Index >= secondCall.Index {
		return nil, errors.New("calls are not sequential")
	}

	// Calls are ordered.
	return []Call{
		firstCall,
		secondCall,
	}, nil
}

func genMatches(current *Node, nodeToGoalMapping map[*Node][]*Goal, goalToDeclsMapping map[*Goal][]Decl) []Decl {
	var all []Decl
	split := ""

	goals := nodeToGoalMapping[current]

	if len(goals) == 0 {
		return nil
	}

	for index, call := range current.Calls {
		wraps := genMatches(call, nodeToGoalMapping, goalToDeclsMapping)

		wrappers := goalToDeclsMapping[goals[index]]
		if len(wrappers) == 0 {
			return nil
		}

		res := wrapDeclSets(wrappers, wraps)

		all = combineDeclSets(all, res, current.Name, split)

		split = call.Name
	}

	return all
}

func combineDeclSets(firstSet []Decl, secondSet []Decl, mustMatch string, mustSplit string) []Decl {
	if len(secondSet) == 0 {
		return nil
	}

	if len(firstSet) == 0 {
		return secondSet
	}

	results := make([]Decl, 0, len(firstSet)*len(secondSet))

	for _, first := range firstSet {
		for _, second := range secondSet {
			combined, err := combineDecls(first, second, mustMatch, mustSplit)
			if err != nil {
				continue
			}

			results = append(results, combined)
		}
	}

	return results
}

// wrapDecl appends the decl tree second onto the end of the linear decl
// first.
func wrapDecl(wrapper Decl, wrapped Decl) Decl {

	if len(wrapper.Calls) == 0 {
		if wrapper.Name != wrapped.Name {
			panic("wrapper name mismatch")
		}

		return wrapped
	}

	lastCall := wrapper.Calls[len(wrapper.Calls)-1]

	return Decl{
		Name:     wrapper.Name,
		Position: wrapper.Position,
		Calls: []Call{
			{
				Name:     lastCall.Name,
				Position: lastCall.Position,
				Index:    lastCall.Index,
				Decl:     wrapDecl(lastCall.Decl, wrapped),
			},
		},
	}
}

func wrapDeclSets(wrappers []Decl, wraps []Decl) []Decl {
	if len(wrappers) == 0 {
		return nil
	}

	if len(wraps) == 0 {
		return wrappers
	}

	results := make([]Decl, 0, len(wrappers)*len(wraps))

	for _, wrapper := range wrappers {
		for _, wrapped := range wraps {
			results = append(results, wrapDecl(wrapper, wrapped))
		}
	}

	return results
}

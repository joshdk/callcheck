// Copyright 2018 Josh Komoroske. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE.txt file.

package cmd

import (
	"fmt"

	"github.com/joshdk/callcheck/config"
	"github.com/joshdk/callcheck/graph"
	"github.com/joshdk/callcheck/policy"
)

func violations(callGraph map[string]graph.FuncDecl, cfg *config.Config) bool {

	var dirty bool

	// Examine each policy
	for _, forbiddenPolicy := range cfg.Forbidden {

		// Find all violations for this policy
		violations := policy.MatchingPaths(callGraph, forbiddenPolicy)

		if len(violations) == 0 {
			continue
		}

		dirty = true

		fmt.Printf("Found %d violations for %s\n", len(violations), forbiddenPolicy.Name)

		for index, violation := range violations {
			if index == 10 {
				fmt.Printf("Violation %d...%d/%d omitted\n", index+1, len(violations), len(violations))
				break
			}

			fmt.Printf("Violation %d/%d\n", index+1, len(violations))
			fmt.Println(violation)
		}

	}

	return dirty
}

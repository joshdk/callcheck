// Copyright 2018 Josh Komoroske. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE.txt file.

package policy

type Policy struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
	Rule        *Node  `yaml:"rule"`
}

type Node struct {
	Name  string  `yaml:"name"`
	Calls []*Node `yaml:"calls"`
}

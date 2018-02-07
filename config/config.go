// Copyright 2018 Josh Komoroske. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE.txt file.

package config

import (
	"github.com/joshdk/callcheck/policy"
)

type Config struct {
	Forbidden []policy.Policy `yaml:"forbid"`
}

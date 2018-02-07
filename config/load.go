// Copyright 2018 Josh Komoroske. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE.txt file.

package config

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

const (
	DefaultFilename = "callcheck.yml"
)

func Load() (*Config, error) {
	return LoadFile(DefaultFilename)
}

func LoadFile(filename string) (*Config, error) {
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var cfg Config

	if err := yaml.Unmarshal(body, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

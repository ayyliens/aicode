package oai

import "github.com/mitranim/gg"

type FunctionCall struct {
	Name      FunctionName `json:"name,omitempty"      yaml:"name,omitempty"      toml:"name,omitempty"`
	Arguments string       `json:"arguments,omitempty" yaml:"arguments,omitempty" toml:"arguments,omitempty"`
}

func (self FunctionCall) IsValid() bool { return gg.IsNotZero(self.Name) }

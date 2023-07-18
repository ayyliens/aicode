package oai

type FunctionDefinition struct {
	Name        string `json:"name,omitempty"        yaml:"name,omitempty"        toml:"name,omitempty"`
	Description string `json:"description,omitempty" yaml:"description,omitempty" toml:"description,omitempty"`
	Parameters  any    `json:"parameters,omitempty"  yaml:"parameters,omitempty"  toml:"parameters,omitempty"`
}

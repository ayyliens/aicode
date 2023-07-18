package oai

type FunctionCall struct {
	Name      string `json:"name,omitempty"      yaml:"name,omitempty"      toml:"name,omitempty"`
	Arguments any    `json:"arguments,omitempty" yaml:"arguments,omitempty" toml:"arguments,omitempty"`
}

func (self FunctionCall) IsValid() bool { return len(self.Name) > 0 }

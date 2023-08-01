package u

import (
	"encoding/json"

	"github.com/mitranim/gg"
	"gopkg.in/yaml.v3"
)

/*
In YAML, this is represented as structured data.

Outside of YAML, this should contain a JSON-encoded string.

When decoding or encoding YAML, this converts between structured data (YAML)
and encoded JSON.
*/
type YamlJsonString string

func (self YamlJsonString) String() string { return string(self) }

func (self YamlJsonString) MarshalYAML() (any, error) {
	if gg.IsZero(self) {
		return nil, nil
	}

	var tar any
	err := json.Unmarshal(gg.ToBytes(self), &tar)
	if err != nil {
		return nil, err
	}
	return tar, nil
}

func (self *YamlJsonString) UnmarshalYAML(src *yaml.Node) error {
	if src == nil {
		gg.PtrClear(self)
		return nil
	}

	var tar any
	err := src.Decode(&tar)
	if err != nil {
		return gg.Wrapf(err, `unable to decode YAML node into %T`, self)
	}

	out, err := json.Marshal(tar)
	if err != nil {
		return gg.Wrapf(err, `unable to JSON-encode value for %T`, self)
	}

	*self = gg.ToText[YamlJsonString](out)
	return nil
}

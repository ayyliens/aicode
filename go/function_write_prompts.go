package main

import (
	"_/go/oai"

	"github.com/mitranim/gg"
)

// Implements `oai.OaiFunction` for writing files.
type FunctionWritePrompts struct {
	Dir oai.ConvDir
}

var _ = oai.OaiFunction(gg.Zero[FunctionWritePrompts]())

func (self FunctionWritePrompts) OaiCall(src string) (_ string) {
	inp := gg.JsonDecodeTo[FunctionWritePromptsInp](src)

	ver := self.Dir.LastVersion()
	for _, file := range inp.Prompts {
		file.WriteTo(ver.NextMinor(), self.Dir.Path)
	}

	return
}

type FunctionWritePromptsInp struct {
	Prompts []oai.Prompt `json:"prompts" desc:"list of prompts with role and contents"`
	// TODO maybe add flag parallel
}

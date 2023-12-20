package main

import (
	"_/go/oai"
	"_/go/u"

	"github.com/mitranim/gg"
)

// Implements `oai.OaiFunction` for writing files.
type FunctionWritePrompts struct {
	Dir oai.ConvDir
}

var _ = oai.OaiFunction(gg.Zero[FunctionWritePrompts]())

func (self FunctionWritePrompts) OaiCall(src string) (_ string) {

	inp := gg.JsonDecodeTo[FunctionWritePromptsInp](src)

	for _, file := range inp.Prompts {
		file.WriteTo(self.Dir.AddMinorVersion(), self.Dir.Path)
	}
	return
}

type FunctionWritePromptsInp struct {
	Prompts []u.Prompt `json:"prompts" desc:"list of prompts with role and contents"`
	// TODO maybe add flag parallel
}

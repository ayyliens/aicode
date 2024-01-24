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

func (self FunctionWritePrompts) Name() oai.FunctionName {
	return `write_prompts`
}

func (self FunctionWritePrompts) OaiCall(ctx u.Ctx, src string) (_ string) {
	inp := gg.JsonDecodeTo[FunctionWritePromptsInp](src)

	ver := self.Dir.LastVersion().AddMinor()
	for _, file := range inp.Prompts {
		file.WriteTo(ver, self.Dir.Path)
		ver = ver.NextMinor()
	}

	return
}

func (self FunctionWritePrompts) Def() oai.FunctionDefinition {
	return oai.FunctionDefinition{
		Name:        string(self.Name()),
		Description: `Provide a list of prompts, with roles and contents`,
		Parameters: map[string]interface{}{
			`type`: `object`,
			`properties`: map[string]interface{}{
				`prompts`: map[string]interface{}{
					`type`:        `array`,
					`description`: `List of prompts`,
					`items`: map[string]interface{}{
						`type`: `object`,
						`properties`: map[string]interface{}{
							`body`: map[string]interface{}{
								`type`:        `string`,
								`description`: `Prompt content.`,
							},
							`role`: map[string]interface{}{
								`type`: `string`,
								`enum`: []string{`user`, `assistant`, `system`},
							},
						},
					},
				},
			},
		},
	}
}

type FunctionWritePromptsInp struct {
	Prompts []oai.Prompt `json:"prompts" desc:"list of prompts with role and contents"`
	// TODO maybe add flag parallel
}

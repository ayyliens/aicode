package main

import (
	"_/go/oai"
	"_/go/u"

	"github.com/mitranim/gg"
)

// Implements `oai.OaiFunction` for writing files.
type FunctionWriteFiles struct {
	Path  string
	Clear bool // Clear output directory before writing.
}

var _ = oai.OaiFunction(gg.Zero[FunctionWriteFiles]())

func (self FunctionWriteFiles) Name() oai.FunctionName {
	return `write_files`
}

func (self FunctionWriteFiles) OaiCall(ctx u.Ctx, src string) (_ string) {
	inp := gg.JsonDecodeTo[FunctionWriteFilesInp](src)

	if self.Clear {
		u.RemoveAllOrSkip(self.Path)
	}

	for _, file := range inp.Files {
		file.WriteTo(self.Path)
	}
	return
}

func (self FunctionWriteFiles) Def() oai.FunctionDefinition {
	return oai.FunctionDefinition{
		Name:        string(self.Name()),
		Description: `Provide a list of files, with file names and contents`,
		Parameters: map[string]interface{}{
			`type`: `object`,
			`properties`: map[string]interface{}{
				`files`: map[string]interface{}{
					`type`:        `array`,
					`description`: `List of files, with file names and contents.`,
					`items`: map[string]interface{}{
						`type`:        `object`,
						`description`: `Individual file, with name and content.`,
						`properties`: map[string]interface{}{
							`body`: map[string]interface{}{
								`type`:        `string`,
								`description`: `File content.`,
							},
							`name`: map[string]interface{}{
								`type`:        `string`,
								`description`: `File name.`,
							},
						},
					},
				},
			},
		},
	}
}

type FunctionWriteFilesInp struct {
	Files []u.File `json:"files" desc:"list of files with names and contents"`
}

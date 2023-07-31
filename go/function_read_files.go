package main

import (
	"_/go/oai"
	"_/go/u"
	"path/filepath"

	"github.com/mitranim/gg"
)

// Implements `oai.OaiFunction` for reading files.
type FunctionReadFiles struct{ Path string }

var _ = oai.OaiFunction(gg.Zero[FunctionReadFiles]())

func (self FunctionReadFiles) OaiCall(src string) (_ string) {
	inp := gg.JsonDecodeTo[FunctionReadFilesInp](src)
	var tar FunctionWriteFilesInp

	for _, path := range inp.Paths {
		var file u.File
		file.Name = path
		file.Validate()
		file.Body = gg.ReadFile[string](filepath.Join(self.Path, path))
		gg.Append(&tar.Files, file)
	}

	if gg.IsZero(tar) {
		return
	}
	return gg.JsonString(tar)
}

type FunctionReadFilesInp struct {
	Paths []string `json:"paths" desc:"list of file paths, relative or absolute"`
}

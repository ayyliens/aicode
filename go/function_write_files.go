package main

import (
	"_/go/oai"

	"github.com/mitranim/gg"
)

/*
Implements `oai.OaiFunction` for writing files.

TODO support option to replace ALL files in target directory.
*/
type FunctionWriteFiles struct{ Path string }

var _ = oai.OaiFunction(gg.Zero[FunctionWriteFiles]())

func (self FunctionWriteFiles) OaiCall(src string) (_ string) {
	inp := gg.JsonDecodeTo[FunctionWriteFilesInput](src)
	for _, file := range inp.Files {
		file.WriteTo(self.Path)
	}
	return
}

type FunctionWriteFilesInput struct {
	Files []File `json:"files" desc:"list of files with names and contents"`
}

package main

import (
	"_/go/oai"

	"github.com/mitranim/gg"
)

/*
Implements `oai.OaiFunction` for writing files.

TODO support option to replace ALL files in target directory.
*/
type FunctionWriteFiles struct {
	Path  string `json:"-"`
	Files []File `json:"files"`
}

var _ = oai.OaiFunction(gg.Zero[FunctionWriteFiles]())

func (self FunctionWriteFiles) OaiCall() (_ string) {
	for _, file := range self.Files {
		file.WriteTo(self.Path)
	}
	return
}

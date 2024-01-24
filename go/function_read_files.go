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

func (self FunctionReadFiles) Name() oai.FunctionName {
	return `read_files`
}

func (self FunctionReadFiles) OaiCall(ctx u.Ctx, src string) (_ string) {
	inp := gg.JsonDecodeTo[FunctionReadFilesInp](src)
	var tar FunctionWriteFilesInp

	// FIXME support wildcards
	// FIXME dedup results

	for _, path := range inp.Paths {
		item := filepath.Join(self.Path, path)
		if path == `*` {
			gg.Append(&tar.Files, LoadDirFiles(self.Path)...)
		} else if gg.DirExists(item) {
			gg.Append(&tar.Files, LoadDirFiles(item)...)
		} else {
			gg.Append(&tar.Files, LoadFile(self.Path, path))
		}
	}

	if gg.IsZero(tar) {
		return
	}
	return gg.JsonString(tar)
}

func LoadDirFiles(path string) (out []u.File) {
	entries := gg.ReadDir(path)
	for _, entry := range entries {
		if entry.IsDir() {
			gg.Append(&out, LoadDirFiles(filepath.Join(path, entry.Name()))...)
		} else {
			gg.Append(&out, LoadFile(path, entry.Name()))
		}
	}
	return
}

func LoadFile(path string, name string) u.File {
	var file u.File
	file.Name = name
	file.Validate()
	file.Body = gg.ReadFile[string](filepath.Join(path, name))
	return file
}

type FunctionReadFilesInp struct {
	Paths []string `json:"paths" desc:"list of file paths, relative or absolute"`
}

func (self FunctionReadFiles) Def() oai.FunctionDefinition {
	return oai.FunctionDefinition{
		Name:        string(self.Name()),
		Description: `Request a list of files by providing a list of file names or paths`,
		Parameters: map[string]interface{}{
			`type`: `object`,
			`properties`: map[string]interface{}{
				`paths`: map[string]interface{}{
					`type`:        `array`,
					`description`: `List of file names or paths.`,
					`items`: map[string]interface{}{
						`type`:        `string`,
						`description`: `File name or path.`,
					},
				},
			},
		},
	}
}

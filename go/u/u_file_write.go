package u

import (
	"path/filepath"

	"github.com/mitranim/gg"
)

/*
Determines how `FileWrite` behaves when `FileWrite.Body` is empty.
Has no effect when `FileWrite.Body` is not empty.
*/
const (
	FileWriteEmptyCreate = `create` // Create target file if missing.
	FileWriteEmptyDelete = `delete` // Delete target file if exists.
	FileWriteEmptyTrunc  = `trunc`  // Truncate target file if exists.
	FileWriteEmptySkip   = `skip`   // Avoid touching target file.
)

/*
TODO consider moving to `gg`.

TODO similar abstraction for various "encode" and "decode" functions such as
`JsonDecodeFileOpt` and `JsonEncodeFileOpt`.
*/
type FileWrite struct {
	Path  string
	Body  []byte
	Mkdir bool
	Empty string
}

func (self FileWrite) Run() {
	if self.Mkdir {
		gg.MkdirAll(filepath.Dir(self.Path))
	}

	if gg.IsTextEmpty(self.Body) {
		switch self.Empty {
		case ``, FileWriteEmptyCreate:
			break

		case FileWriteEmptyDelete:
			RemoveFileOrDirOpt(self.Path)

		case FileWriteEmptyTrunc:
			if !gg.FileExists(self.Path) {
				return
			}

		case FileWriteEmptySkip:
			return

		default:
			panic(gg.Errf(`unknown FileWrite.Empty: %q`, self.Empty))
		}
	}

	gg.WriteFile(self.Path, self.Body)
}

func WriteFile[A gg.Text](path string, src A) {
	FileWrite{Path: path, Body: gg.ToBytes(src)}.Run()
}

func WriteFileOpt[A gg.Text](path string, src A) {
	FileWrite{Path: path, Body: gg.ToBytes(src), Empty: FileWriteEmptyTrunc}.Run()
}

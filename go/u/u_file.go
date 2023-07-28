package u

import (
	"path/filepath"

	"github.com/mitranim/gg"
)

type File struct {
	Name string `json:"name" desc:"file name (must be local, without directory)"`
	Body string `json:"body" desc:"file content"`
}

func (self File) Validate() {
	if gg.IsZero(self.Name) {
		panic(gg.Errv(`missing file name`))
	}
	if !filepath.IsLocal(self.Name) {
		panic(gg.Errf(`unexpected non-local file name %q`, self.Name))
	}
}

func (self File) WriteTo(out string) {
	self.Validate()
	WriteFileRec(filepath.Join(out, self.Name), self.Body)
}

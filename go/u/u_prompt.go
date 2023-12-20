package u

import (
	"path/filepath"

	"github.com/mitranim/gg"
)

type Prompt struct {
	Role string `json:"role" desc:"role (system, function, user)"`
	Body string `json:"body" desc:"file content"`
}

func (self Prompt) Validate() {
	if gg.IsZero(self.Role) {
		panic(gg.Errv(`missing role`))
	}
}

func (self Prompt) WriteTo(ver Version, out string) {
	self.Validate()
	// TODO support yaml
	WriteFileRec(filepath.Join(out, ver.String()+`_`+self.Role+`_message.md`), self.Body)
}

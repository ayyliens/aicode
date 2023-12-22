package oai

import (
	"_/go/u"
	"path/filepath"

	"github.com/mitranim/gg"
)

type Prompt struct {
	Role ChatMessageRole `json:"role" desc:"role (system, function, user)"`
	Body string          `json:"body" desc:"file content"`
}

func (self Prompt) Validate() {
	if gg.IsZero(self.Role) {
		panic(gg.Errv(`missing role`))
	}
}

func (self Prompt) WriteTo(ver u.Version, out string) {
	self.Validate()

	// TODO support yaml
	fileName := VersionedFileName{
		Version: ver,
		Role:    self.Role,
		Type:    VersionedFileTypeMessage,
		Ext:     ChatCompletionMessageDefaultExtForText,
	}.String()

	u.WriteFileRec(filepath.Join(out, fileName), self.Body)
}

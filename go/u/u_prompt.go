package u

import (
	"_/go/oai"
	"path/filepath"

	"github.com/mitranim/gg"
)

type Prompt struct {
	Role oai.ChatMessageRole `json:"role" desc:"role (system, function, user)"`
	Body string              `json:"body" desc:"file content"`
}

func (self Prompt) Validate() {
	if gg.IsZero(self.Role) {
		panic(gg.Errv(`missing role`))
	}
}

func (self Prompt) WriteTo(ver Version, out string) {
	self.Validate()
	// TODO support yaml
	fileName := oai.VersionedFileName{
		Version: ver,
		Role:    self.Role,
		Type:    oai.VersionedFileTypeMessage,
		Ext:     oai.ChatCompletionMessageDefaultExtForText,
	}.String()
	WriteFileRec(filepath.Join(out, fileName), self.Body)
}

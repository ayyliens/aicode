package main

import (
	"_/go/u"
	"path/filepath"

	"github.com/mitranim/gg"
)

type File struct {
	Name string `json:"name"`
	Body string `json:"body"`
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

	u.FileWrite{
		Path:  filepath.Join(out, self.Name),
		Body:  gg.ToBytes(self.Body),
		Mkdir: true,
	}.Run()
}

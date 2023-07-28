package main

import (
	"_/go/u"
	"fmt"

	"github.com/mitranim/cmd"
	"github.com/mitranim/gg"
	"github.com/mitranim/gg/grepr"
	"github.com/mitranim/oas"
)

type CmdPrintOasDoc[_ any] struct {
	OutPath string `flag:"--out-path" desc:"path for output file"`
}

func (self CmdPrintOasDoc[_]) RunCli() {
	gg.FlagParse(cmd.Args(), &self)
	self.Run()
}

func (self CmdPrintOasDoc[A]) Run() {
	typ := gg.Type[A]()

	var doc oas.Doc
	doc.TypeSchema(typ)

	src := doc

	if gg.IsZero(self.OutPath) {
		grepr.Println(src)
		fmt.Println(u.JsonEncodePretty[string](src))
	} else {
		u.PolyEncodeFileOpt(self.OutPath, src)
	}
}

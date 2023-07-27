package main

import (
	"_/go/u"
	"log"

	"github.com/mitranim/cmd"
	"github.com/mitranim/gg"
)

type CmdForkDir struct {
	Path string `flag:"--path" desc:"source path"`
}

func (self CmdForkDir) RunCli() {
	gg.FlagParse(cmd.Args(), &self)
	self.Run()
}

func (self CmdForkDir) Run() {
	src := self.Path

	if gg.IsZero(src) {
		panic(gg.Errf(`missing path: "--path"`))
	}

	tar := u.IndexedDirForkPath(src)

	log.Printf(`copying %q to %q`, src, tar)
	u.CopyDirRec(src, tar)
}

package main

import (
	"_/go/oai"
	"context"

	"github.com/mitranim/cmd"
	"github.com/mitranim/gg"
)

/*
Example usage:

	make go.run   run='oai_conv_file --verb --path local/conv.json'
	make go.run.w run='oai_conv_file --verb --path local/conv.json --watch --init'
*/
type CmdOaiConvFile struct{ oai.OaiClientConvFile }

func (self CmdOaiConvFile) RunCli() {
	gg.FlagParse(cmd.Args(), &self)
	self.Run()
}

func (self CmdOaiConvFile) Run() {
	self.OaiClientConvFile.Run(context.Background())
}

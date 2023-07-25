package main

import (
	"_/go/oai"
	"context"

	"github.com/mitranim/cmd"
	"github.com/mitranim/gg"
)

/*
Example usage:

	make go.run run='cmd_conv_file --path local/conv.json'
	make go.run.w run='cmd_conv_file --path local/conv.json --watch --init'
*/
type CmdOaiConvFile struct {
	CmdOaiCommon
	ApiKey string
}

func (self CmdOaiConvFile) RunCli() {
	gg.FlagParse(cmd.Args(), &self)
	self.Run()
}

func (self CmdOaiConvFile) Run() {
	if gg.IsZero(self.Path) {
		panic(gg.Errf(`missing path: "--path"`))
	}

	ctx := context.Background()

	var cli oai.OaiClientConvFile
	cli.ApiKey = self.ApiKey
	cli.Path = self.Path
	cli.Verb = true

	if self.Watch {
		cli.Init = self.Init
		cli.Watch(ctx)
	} else {
		cli.Run(ctx)
	}
}

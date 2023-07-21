package main

import (
	"_/go/oai"
	"context"

	"github.com/mitranim/cmd"
	"github.com/mitranim/gg"
)

/*
Example usage:

	make go.run run='cmd_conv_dir --path local/conv'
	make go.run.w run='cmd_conv_dir --path local/conv --watch --init'
*/
type CmdOaiConvDir struct {
	Path string `flag:"--path" desc:"directory path (required)"`
	CmdOaiCommon
	ApiKey string
}

func (self CmdOaiConvDir) RunCli() {
	gg.FlagParse(cmd.Args(), &self)
	self.Run()
}

func (self CmdOaiConvDir) Run() {
	if gg.IsZero(self.Path) {
		panic(gg.Errf(`missing path: "--path"`))
	}

	ctx := context.Background()

	var cli oai.OaiClientConvDir
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

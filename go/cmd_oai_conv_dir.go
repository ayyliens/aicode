package main

import (
	"_/go/oai"
	"context"

	"github.com/mitranim/cmd"
	"github.com/mitranim/gg"
)

/*
Examples:

	make go.run run='cmd_conv_dir --path=local/conv'
	make go.run.w run='cmd_conv_dir --path=local/conv --watch --init'
	make go.run.w run='cmd_conv_dir --path=local/conv --out-path=local/conv/files --watch --funcs'
*/
type CmdOaiConvDir struct {
	CmdOaiCommon
	ApiKey  string
	Path    string `flag:"--path"     desc:"directory path (required)"`
	OutPath string `flag:"--out-path" desc:"directory path for output files"`
	Funcs   bool   `flag:"--funcs"    desc:"automatically run registered functions"`
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

	if self.Funcs && gg.IsNotZero(self.OutPath) {
		cli.Functions.Add(`write_files`, &FunctionWriteFiles{
			Path: self.OutPath,
		})
	}

	if self.Watch {
		cli.Init = self.Init
		cli.Watch(ctx)
	} else {
		cli.Run(ctx)
	}
}

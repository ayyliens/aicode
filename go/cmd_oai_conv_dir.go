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

Be cautious: files in target directory may be overwritten with no recovery.
*/
type CmdOaiConvDir struct {
	CmdOaiCommon
	ApiKey  string
	OutPath string `flag:"--out-path" desc:"directory path for output files"`
	Funcs   bool   `flag:"--funcs"    desc:"automatically run registered functions"`
	Trunc   bool   `flag:"--trunc"    desc:"support conversation truncation"`
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

	// TODO consider deduping fields between lib type and CLI type.
	// The duplication exists because it's not traditional for library
	// utility types to define CLI "flag" field annotations.
	var cli oai.OaiClientConvDir
	cli.ApiKey = self.ApiKey
	cli.Path = self.Path
	cli.Trunc = self.Trunc
	cli.Verb = true

	if self.Funcs {
		cli.Functions.Add(`get_current_weather`, &FunctionGetCurrentWeather{})

		if gg.IsNotZero(self.OutPath) {
			cli.Functions.Add(`write_files`, &FunctionWriteFiles{Path: self.OutPath})
		}
	}

	if self.Watch {
		cli.Init = self.Init
		cli.Watch(ctx)
	} else {
		cli.Run(ctx)
	}
}

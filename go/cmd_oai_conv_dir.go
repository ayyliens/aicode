package main

import (
	"_/go/oai"
	"context"

	"github.com/mitranim/cmd"
	"github.com/mitranim/gg"
)

/*
Examples:

	make go.run run='oai_conv_dir --path=local/conv'
	make go.run.w run='oai_conv_dir --path=local/conv --watch --init'
	make go.run.w run='oai_conv_dir --path=local/conv --watch --funcs --out-path=local/conv/files'
	make go.run.w run='oai_conv_dir --path=local/conv --watch --funcs --out-path=local/conv/files --trunc --fork'

Be cautious: files in target directory may be overwritten with no recovery.
*/
type CmdOaiConvDir struct {
	CmdOaiCommon
	ApiKey  string
	SrcPath string `flag:"--src-path" desc:"directory path for reading source files"`
	OutPath string `flag:"--out-path" desc:"directory path for writing output files"`
	Funcs   bool   `flag:"--funcs"    desc:"automatically run registered functions"`
	Trunc   bool   `flag:"--trunc"    desc:"support conversation truncation in watch mode"`
	Fork    bool   `flag:"--fork"     desc:"support conversation forking in watch mode (best with --trunc)"`
	Dry     bool   `flag:"--dry"      desc:"dry run: no request to external API"`
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
	cli.Fork = self.Fork
	cli.Dry = self.Dry
	cli.Verb = true

	if self.Funcs {
		cli.Functions.Add(`get_current_weather`, FunctionGetCurrentWeather{})

		if gg.IsNotZero(self.SrcPath) {
			cli.Functions.Add(`read_files`, FunctionReadFiles{Path: self.SrcPath})
		}

		if gg.IsNotZero(self.OutPath) {
			cli.Functions.Add(`write_files`, FunctionWriteFiles{Path: self.OutPath})
		}
	}

	if self.Watch {
		cli.Init = self.Init
		cli.Watch(ctx)
	} else {
		cli.Run(ctx)
	}
}

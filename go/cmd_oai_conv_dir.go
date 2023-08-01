package main

import (
	"_/go/oai"
	"_/go/u"
	"context"

	"github.com/mitranim/cmd"
	"github.com/mitranim/gg"
)

/*
Examples:

	make go.run   run='oai_conv_dir --verb --path=local/conv'
	make go.run.w run='oai_conv_dir --verb --path=local/conv --watch --init'
	make go.run.w run='oai_conv_dir --verb --path=local/conv --watch --funcs --out-path=local/conv/files'
	make go.run.w run='oai_conv_dir --verb --path=local/conv --watch --funcs --out-path=local/conv/files --trunc --fork'

Be cautious: files in target directory may be overwritten with no recovery.
*/
type CmdOaiConvDir struct {
	oai.ClientConvDir
	SrcPath string `flag:"--src-path" desc:"directory path for reading source files"`
	OutPath string `flag:"--out-path" desc:"directory path for writing output files"`
	Funcs   bool   `flag:"--funcs"    desc:"automatically run supported functions"`
}

func (self CmdOaiConvDir) RunCli() {
	gg.FlagParse(cmd.Args(), &self)
	self.Run()
}

func (self CmdOaiConvDir) Run() {
	if gg.IsZero(self.Path) {
		panic(gg.Errf(`missing path: "--path"`))
	}

	if self.Funcs {
		self.Functions.Add(`get_current_weather`, FunctionGetCurrentWeather{})

		if gg.IsNotZero(self.SrcPath) {
			self.Functions.Add(`read_files`, FunctionReadFiles{Path: self.SrcPath})
		}

		if gg.IsNotZero(self.OutPath) {
			self.Functions.Add(`write_files`, FunctionWriteFiles{Path: self.OutPath})
		}
	}

	self.Ignore = u.AdjoinCompact(self.Ignore, self.SrcPath, self.OutPath)
	self.ClientConvDir.Run(context.Background())
}

package main

import (
	"_/go/oai"
	"context"
	"log"
	"os"

	"github.com/mitranim/cmd"
	"github.com/mitranim/gg"
)

func main() {
	log.SetFlags(0)

	defer gg.Fatal()

	cmd.Map{
		`oai_watch_file`: CmdOaiWatchFile,
		`oai_watch_dir`:  CmdOaiWatchDir,
	}.Get()()
}

func CmdOaiWatchFile() {
	var cli oai.OaiClientConvFile
	cli.ApiKey = os.Getenv(`OPENAI_API_KEY`)
	cli.Path = `local/conv.json`
	cli.Verb = true
	cli.Watch(context.Background())
}

func CmdOaiWatchDir() {
	var cli oai.OaiClientConvDir
	cli.ApiKey = os.Getenv(`OPENAI_API_KEY`)
	cli.Path = `local/conv`
	cli.Verb = true
	cli.Watch(context.Background())
}

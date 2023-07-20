package main

import (
	"_/go/oai"
	"context"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/joho/godotenv"
	"github.com/mitranim/cmd"
	"github.com/mitranim/gg"
)

func main() {
	log.SetFlags(0)
	defer gg.Fatal()
	loadEnvFiles()

	cmd.Map{
		`oai_watch_file`: CmdOaiWatchFile,
		`oai_watch_dir`:  CmdOaiWatchDir{}.RunCli,
	}.Get()()
}

func CmdOaiWatchFile() {
	var cli oai.OaiClientConvFile
	cli.ApiKey = os.Getenv(`OPENAI_API_KEY`)
	cli.Path = `local/conv.json`
	cli.Verb = true
	cli.Watch(context.Background())
}

type CmdOaiWatchDir struct {
	Path string `flag:"-p"`
}

func (self CmdOaiWatchDir) RunCli() {
	gg.FlagParse(cmd.Args(), &self)
	self.Run()
}

func (self CmdOaiWatchDir) Run() {
	if gg.IsZero(self.Path) {
		panic(gg.Errf(`missing path: "-p"`))
	}

	var cli oai.OaiClientConvDir
	cli.ApiKey = os.Getenv(`OPENAI_API_KEY`)
	cli.Path = self.Path
	cli.Verb = true
	cli.Watch(context.Background())
}

func loadEnvFiles() {
	for _, base := range gg.Reversed(strings.Split(os.Getenv(`CONF`), `,`)) {
		gg.Try(godotenv.Load(filepath.Join(base, `.env.properties`)))
	}
}

package main

import (
	"_/go/u"
	"log"
	"os"

	"github.com/mitranim/cmd"
	"github.com/mitranim/gg"
)

func main() {
	log.SetFlags(0)
	defer gg.Fatal()
	u.LoadEnvFiles()

	key := os.Getenv(`OPENAI_API_KEY`)

	cmd.Map{
		`oai_conv_file`: CmdOaiConvFile{ApiKey: key}.RunCli,
		`oai_conv_dir`:  CmdOaiConvDir{ApiKey: key}.RunCli,
	}.Get()()
}

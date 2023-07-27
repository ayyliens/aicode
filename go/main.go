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

		// May be helpful for generating JSON schemas for request templates.
		// Substitute the input type as needed.
		`print_oas_doc`: CmdPrintOasDoc[FunctionWriteFilesInput]{
			OutPath: `local/oas_doc.yaml`,
		}.RunCli,
	}.Get()()
}

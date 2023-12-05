package main

import (
	"_/go/oai"
	"log"
	"os"
	"path/filepath"

	"github.com/joeshaw/envdecode"
	"github.com/joho/godotenv"
	"github.com/mitranim/cmd"
	"github.com/mitranim/gg"
)

func main() {
	log.SetFlags(0)
	defer gg.Fatal()

	var conf Conf
	conf.Init()

	client := conf.OaiClient()

	cmd.Map{
		`oai_conv_file`: gg.With(func(tar *CmdOaiConvFile) {
			tar.Client = client
		}).RunCli,

		`oai_conv_dir`: gg.With(func(tar *CmdOaiConvDir) {
			tar.Client = client
			tar.Model = conf.OpenAiModel
		}).RunCli,

		`fork_dir`: CmdForkDir{}.RunCli,

		// May be helpful for generating JSON schemas for request templates.
		// Substitute the input type as needed.
		`print_oas_doc`: CmdPrintOasDoc[FunctionWriteFilesInp]{
			OutPath: `local/oas_doc.yaml`,
		}.RunCli,
	}.Get()()
}

type Conf struct {
	OpenAiApiKey string `env:"OPEN_AI_API_KEY"`
	OpenAiMock   bool   `env:"OPEN_AI_MOCK"`
	OpenAiModel  string `env:"OPEN_AI_MODEL"`
}

func (self *Conf) Init() {
	LoadEnvFileOpt(os.Getenv(`CONF`))
	LoadEnvFileOpt(`.`)
	gg.Try(envdecode.StrictDecode(self))
}

func LoadEnvFileOpt(path string) {
	if gg.IsNotZero(path) {
		gg.Try(godotenv.Load(filepath.Join(path, `.env.properties`)))
	}
}

func (self Conf) OaiClient() oai.Client {
	if self.OpenAiMock {
		return gg.Zero[oai.MockClient]()
	}

	if gg.IsNotZero(self.OpenAiApiKey) {
		var tar oai.HttpClient
		tar.ApiKey = self.OpenAiApiKey
		return tar
	}

	panic(gg.Errv(`unable to make OpenAI client: missing API key, and mocks not enabled`))
}

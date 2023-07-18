package oai_test

import (
	"_/go/oai"
	"_/go/u"
	"context"
	"testing"

	"github.com/mitranim/gg"
	"github.com/mitranim/gg/grepr"
	"github.com/mitranim/gg/gtest"
)

func Test_chat_completion(t *testing.T) {
	t.Skip()
	t.Fail()
	defer gtest.Catch(t)

	var cli oai.OaiClient
	cli.Init()

	var req oai.ChatCompletionRequest
	req.Default()
	req.Messages = u.ParseJsonLines[oai.ChatCompletionMessage](gg.ReadFile[string](`testdata/conv_1_to_api_compatible.json`))

	grepr.Prn(`req.Messages:`, req.Messages)

	gg.Append(&req.Messages, oai.ChatCompletionMessage{
		Role:    oai.ChatMessageRoleUser,
		Content: `Summarize the conversation so far.`,
	})

	ctx := context.Background()
	res := cli.ChatCompletion(ctx, req)
	grepr.Prn(`res:`, res)
}

func Test_conv_conversion(t *testing.T) {
	t.Skip()
	defer gtest.Catch(t)

	src := gg.ReadFile[string](`testdata/conv_1.json`)
	mid := oai.OaiSiteMsgs(u.ParseJsonLines[oai.OaiSiteMsg](src))
	tar := mid.ChatCompletionMessages()

	u.WriteJsonLines(`testdata/conv_1_to_api_compatible.json`, tar)
}

func Test_parsing_conv_json(t *testing.T) {
	// t.Fail()
	defer gtest.Catch(t)

	src := gg.ReadFile[string](`testdata/conv_0.json`)
	tar := u.ParseJsonLines[oai.ChatCompletionMessage](src)

	grepr.Prn(`tar:`, tar)
}

func Test_conv_combine(t *testing.T) {
	t.Skip()
	defer gtest.Catch(t)

	src := gg.ReadFile[string](`testdata/conv_1.json`)
	tar := oai.OaiSiteMsgs(u.ParseJsonLines[oai.OaiSiteMsg](src))

	gg.WriteFile(`testdata/conv_1_combined.txt`, tar.ConcatText())
}

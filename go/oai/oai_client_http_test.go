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

func Test_HttpClient_ChatCompletionRequest(t *testing.T) {
	t.Skip(`enable on demand when needed; requires API key`)
	t.Fail()
	defer gtest.Catch(t)

	var req oai.ChatCompletionRequest
	req.Default()
	req.Messages = u.ParseJsonLines[oai.ChatCompletionMessage](gg.ReadFile[string](`testdata/conv_hello_world.json`))

	grepr.Prn(`req.Messages:`, req.Messages)

	gg.Append(&req.Messages, oai.ChatCompletionMessage{
		Role:    oai.ChatMessageRoleUser,
		Content: `Summarize the conversation so far.`,
	})

	var cli oai.HttpClient
	cli.ApiKey = ``

	ctx := context.Background()
	res := cli.ChatCompletionResponse(ctx, req)
	grepr.Prn(`res:`, res)
}

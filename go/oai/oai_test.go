package oai_test

import (
	"_/go/oai"
	"_/go/u"
	"context"
	"strings"
	"testing"

	"github.com/mitranim/gg"
	"github.com/mitranim/gg/grepr"
	"github.com/mitranim/gg/gtest"
)

func Test_chat_completion(t *testing.T) {
	t.Skip()
	t.Fail()
	defer gtest.Catch(t)

	var req oai.ChatCompletionRequest
	req.Default()
	req.Messages = u.ParseJsonLines[oai.ChatCompletionMessage](gg.ReadFile[string](`testdata/conv_1_to_api_compatible.json`))

	grepr.Prn(`req.Messages:`, req.Messages)

	gg.Append(&req.Messages, oai.ChatCompletionMessage{
		Role:    oai.ChatMessageRoleUser,
		Content: `Summarize the conversation so far.`,
	})

	var cli oai.OaiClient
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

func Test_OaiSiteMsgs_ChatCompletionMessages(t *testing.T) {
	defer gtest.Catch(t)

	gtest.Equal(
		testOaiSiteMsgJsonLines().ChatCompletionMessages(),
		[]oai.ChatCompletionMessage{
			{Role: `user`, Content: `provide response 0`},
			{Role: `assistant`, Content: `response 0 provided`},
			{Role: `user`, Content: `provide response 1`},
			{Role: `assistant`, Content: `response 1 provided`},
		},
	)
}

func testOaiSiteMsgJsonLines() oai.OaiSiteMsgs {
	return u.ParseJsonLines[oai.OaiSiteMsg](
		gg.ReadFile[string](`testdata/site_msgs_0.json`),
	)
}

func Test_OaiSiteMsgs_ConcatText(t *testing.T) {
	defer gtest.Catch(t)

	gtest.Eq(
		testOaiSiteMsgJsonLines().ConcatText(),
		strings.TrimSpace(`
provide response 0

response 0 provided

provide response 1

response 1 provided
`),
	)
}

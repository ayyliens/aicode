package oai_test

import (
	"_/go/oai"
	"_/go/u"
	"strings"
	"testing"

	"github.com/mitranim/gg"
	"github.com/mitranim/gg/gtest"
)

func Test_SiteMsgs_ChatCompletionMessages(t *testing.T) {
	defer gtest.Catch(t)

	gtest.Equal(
		testSiteMsgJsonLines().ChatCompletionMessages(),
		[]oai.ChatCompletionMessage{
			{Role: `user`, Content: `provide response 0`},
			{Role: `assistant`, Content: `response 0 provided`},
			{Role: `user`, Content: `provide response 1`},
			{Role: `assistant`, Content: `response 1 provided`},
		},
	)
}

func testSiteMsgJsonLines() oai.SiteMsgs {
	return u.ParseJsonLines[oai.SiteMsg](
		gg.ReadFile[string](`testdata/site_msgs.json`),
	)
}

func Test_SiteMsgs_ConcatText(t *testing.T) {
	defer gtest.Catch(t)

	gtest.Eq(
		testSiteMsgJsonLines().ConcatText(),
		strings.TrimSpace(`
provide response 0

response 0 provided

provide response 1

response 1 provided
`),
	)
}

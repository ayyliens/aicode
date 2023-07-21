package oai

import "github.com/mitranim/gg"

// Represents the format used internally by the chat UI on the OpenAI website.
type OaiSiteMsg struct {
	Author  OaiSiteMsgAuthor  `json:"author" db:"author"`
	Content OaiSiteMsgContent `json:"content" db:"content"`
}

func (self OaiSiteMsg) GetText() string { return self.Content.GetText() }

func (self OaiSiteMsg) ChatCompletionMessage() (out ChatCompletionMessage) {
	src := self.GetText()
	if gg.IsZero(src) {
		return
	}

	out.Role = self.Author.Role
	out.Content = self.GetText()
	return
}

type OaiSiteMsgAuthor struct {
	Role ChatMessageRole `json:"role" db:"role"`
}

type OaiSiteMsgContent struct {
	ContentType string   `json:"content_type" db:"content_type"`
	Parts       []string `json:"parts" db:"parts"`
}

func (self OaiSiteMsgContent) GetText() string {
	return gg.JoinSpacedOpt(self.Parts...)
}

package oai

import "github.com/mitranim/gg"

// Represents the format used internally by the chat UI on the OpenAI website.
type SiteMsg struct {
	Author  SiteMsgAuthor  `json:"author" db:"author"`
	Content SiteMsgContent `json:"content" db:"content"`
}

func (self SiteMsg) GetText() string { return self.Content.GetText() }

func (self SiteMsg) ChatCompletionMessage() (out ChatCompletionMessage) {
	src := self.GetText()
	if gg.IsZero(src) {
		return
	}

	out.Role = self.Author.Role
	out.Content = self.GetText()
	return
}

type SiteMsgAuthor struct {
	Role ChatMessageRole `json:"role" db:"role"`
}

type SiteMsgContent struct {
	ContentType string   `json:"content_type" db:"content_type"`
	Parts       []string `json:"parts" db:"parts"`
}

func (self SiteMsgContent) GetText() string {
	return gg.JoinSpacedOpt(self.Parts...)
}

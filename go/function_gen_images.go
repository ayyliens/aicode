package main

import (
	"_/go/oai"
	"_/go/u"
	"encoding/base64"
	"path/filepath"

	"github.com/mitranim/gg"
)

// Implements `oai.OaiFunction` for image generation.
type FunctionGenImages struct {
	Client oai.Client
	Path   string
	Clear  bool // Clear output directory before writing.
}

var _ = oai.OaiFunction(gg.Zero[FunctionGenImages]())

func (self FunctionGenImages) Name() oai.FunctionName {
	return `gen_images`
}

func (self FunctionGenImages) OaiCall(ctx u.Ctx, src string) (_ string) {
	inp := gg.JsonDecodeTo[FunctionGenImagesInp](src)

	if self.Clear {
		u.RemoveAllOrSkip(self.Path)
	}

	for _, prompt := range inp.Prompts {
		res := self.Client.ImageGenerationResponse(ctx, oai.ImageGenerationRequest{
			Prompt: prompt.Prompt,
			N:      1,
			// FIXME control image size
			//Size: oai.ImageGenerationSize1024x1024,
			ResponseFormat: oai.ImageGenerationFormatB64JSON,
		})
		for _, data := range res.Data {
			decoded := gg.Try1(base64.StdEncoding.DecodeString(data.B64Json))
			u.WriteFileRec(filepath.Join(self.Path, prompt.Name+`.png`), decoded)
		}
	}
	return
}

func (self FunctionGenImages) Def() oai.FunctionDefinition {
	return oai.FunctionDefinition{
		Name:        string(self.Name()),
		Description: `Provide a list of image prompts`,
		Parameters: map[string]interface{}{
			`type`: `object`,
			`properties`: map[string]interface{}{
				`prompts`: map[string]interface{}{
					`type`:        `array`,
					`description`: `List of prompts to generate image`,
					`items`: map[string]interface{}{
						`type`:        `object`,
						`description`: `Prompt for image generation, with image file name`,
						`properties`: map[string]interface{}{
							`prompt`: map[string]interface{}{
								`type`:        `string`,
								`description`: `Image prompt`,
							},
							`name`: map[string]interface{}{
								`type`:        `string`,
								`description`: `File name (must be local, without directory).`,
							},
						},
					},
				},
			},
		},
	}
}

type FunctionGenImagesInp struct {
	Prompts []ImagePrompt `json:"prompts" desc:"list of image prompts"`
}

type ImagePrompt struct {
	Name   string `json:"name" desc:"file name (must be local, without directory)"`
	Prompt string `json:"prompt" desc:"Prompt for image generation"`
}

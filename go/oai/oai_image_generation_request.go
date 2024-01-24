package oai

type ImageGenerationSize string

const (
	ImageGenerationSize256x256   = ImageGenerationSize("256x256")
	ImageGenerationSize512x512   = ImageGenerationSize("512x512")
	ImageGenerationSize1024x1024 = ImageGenerationSize("1024x1024")
)

type ImageGenerationFormat string

const (
	ImageGenerationFormatURL     = ImageGenerationFormat("url")
	ImageGenerationFormatB64JSON = ImageGenerationFormat("b64_json")
)

type ImageGenerationRequest struct {
	Prompt         string                `json:"prompt,omitempty"`
	N              int                   `json:"n,omitempty"`
	Size           ImageGenerationSize   `json:"size,omitempty"`
	ResponseFormat ImageGenerationFormat `json:"response_format,omitempty"`
	User           string                `json:"user,omitempty"`
}

type ImageGenerationResponse struct {
	Created int64                 `json:"created,omitempty"`
	Data    []ImageGenerationData `json:"data,omitempty"`
}

type ImageGenerationData struct {
	Url     string `json:"url,omitempty"`
	B64Json string `json:"b64_json,omitempty"`
}

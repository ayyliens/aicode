package main

import (
	"_/go/oai"
	"_/go/u"

	"github.com/mitranim/gg"
)

type FunctionGetCurrentWeather struct{}

var _ = oai.OaiFunction(gg.Zero[FunctionGetCurrentWeather]())

func (self FunctionGetCurrentWeather) Name() oai.FunctionName {
	return `get_current_weather`
}

func (FunctionGetCurrentWeather) OaiCall(ctx u.Ctx, src string) string {
	inp := gg.JsonDecodeTo[FunctionGetCurrentWeatherInp](src)

	return gg.JsonString(FunctionGetWeatherOut{
		Temperature: 23,
		Unit:        gg.Or(inp.Unit, `celsius`),
		Description: `sunny`,
	})
}

func (self FunctionGetCurrentWeather) Def() oai.FunctionDefinition {
	return oai.FunctionDefinition{
		Name:        string(self.Name()),
		Description: `Get the current weather in a given location`,
		Parameters: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"location": map[string]interface{}{
					"type":        "string",
					"description": "The city and state, e.g. San Francisco, CA",
				},
				"unit": map[string]interface{}{
					"type": "string",
					"enum": []string{"celsius", "fahrenheit"},
				},
			},
			"required": []string{"location"},
		},
	}
}

type FunctionGetCurrentWeatherInp struct {
	Location string `json:"location"`
	Unit     string `json:"unit"`
}

type FunctionGetWeatherOut struct {
	Temperature float64 `json:"temperature"`
	Unit        string  `json:"unit"`
	Description string  `json:"description"`
}

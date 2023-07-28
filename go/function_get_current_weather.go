package main

import (
	"_/go/oai"

	"github.com/mitranim/gg"
)

type FunctionGetCurrentWeather struct{}

var _ = oai.OaiFunction(gg.Zero[FunctionGetCurrentWeather]())

func (FunctionGetCurrentWeather) OaiCall(src string) string {
	inp := gg.JsonDecodeTo[FunctionGetCurrentWeatherInp](src)

	return gg.JsonString(FunctionGetWeatherOut{
		Temperature: 23,
		Unit:        gg.Or(inp.Unit, `celsius`),
		Description: `sunny`,
	})
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

package main

import (
	"_/go/oai"

	"github.com/mitranim/gg"
)

type FunctionGetCurrentWeather struct {
	Location string `json:"location"`
	Unit     string `json:"unit"`
}

var _ = oai.OaiFunction(gg.Zero[FunctionGetCurrentWeather]())

func (self FunctionGetCurrentWeather) OaiCall() string {
	return gg.JsonString(FunctionResponseGetWeather{
		Temperature: 23,
		Unit:        gg.Or(self.Unit, `celsius`),
		Description: `sunny`,
	})
}

type FunctionResponseGetWeather struct {
	Temperature float64 `json:"temperature"`
	Unit        string  `json:"unit"`
	Description string  `json:"description"`
}

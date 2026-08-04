package costable

import (
	_ "embed"
	"gopkg.in/yaml.v3"
)

type ModelPricing struct {
	InputPerMillion  float64 `yaml:"input"`
	OutputPerMillion float64 `yaml:"output"`
}

type PricingData struct {
	Models map[string]ModelPricing `yaml:"models"`
}

//go:embed pricing.yaml
var pricingFile []byte

var data PricingData

func init() {
	err := yaml.Unmarshal(pricingFile, &data)
	if err != nil {
		panic(err)
	}
}

func ValidateModel(model string) bool {
	_, ok := data.Models[model]

	return ok
}

func CalculateCost(model string, inputTokens, outputTokens int) float64 {

	modelPricing, ok := data.Models[model]

	if !ok {
		return 0
	}

	return (float64(inputTokens) * modelPricing.InputPerMillion / 1_000_000) +
		(float64(outputTokens) * modelPricing.OutputPerMillion / 1_000_000)
}

package provider

import (
	"fmt"
	"os"
)

type ExtractedData struct {
	Model        string
	SystemPrompt string
	UserPrompt   string
	Response     string
	Tokens       int
	Cost         int
}

type Extractor interface {
	ExtractRequest(body []byte, path string) (model, systemPrompt, userPrompt string)
	ExtractResponse(body []byte) (response string, tokens int)
}

func GetExtractor(host string) Extractor {
	switch hostProvider[host] {
	case "openai":
		return &OpenAIExtractor{}
	case "gemini":
		return &GeminiExtractor{}
	case "anthropic":
		return &AnthropicExtractor{}
	default:
		return nil
	}
}

var modelProvider = map[string]string{
	"gpt-4o":           "openai",
	"gpt-4o-mini":      "openai",
	"claude-3-sonnet":  "anthropic",
	"claude-3-haiku":   "anthropic",
	"gemini-2.5-flash": "gemini",
}

var hostProvider = map[string]string{
	"api.openai.com":                    "openai",
	"api.anthropic.com":                 "anthropic",
	"generativelanguage.googleapis.com": "gemini",
}

func GetAPIKey(host string) (string, error) {
	switch host {
	case "api.openai.com":
		return os.Getenv("OPENAI_API_KEY"), nil
	case "api.anthropic.com":
		return os.Getenv("ANTHROPIC_API_KEY"), nil
	case "generativelanguage.googleapis.com":
		return os.Getenv("GEMINI_API_KEY"), nil
	default:
		return "", fmt.Errorf("unknown provider: %s", host)
	}
}

func ValidateModel(modelFlag, traceHost string) error {
	if hostProvider[traceHost] != modelProvider[modelFlag] {
		return fmt.Errorf("cannot replay: trace was sent to %s but model %s belongs to %s", traceHost, modelFlag, modelProvider[modelFlag])
	}
	return nil
}

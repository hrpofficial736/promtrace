package provider

import (
	"fmt"
	"os"
)

var modelProvider = map[string]string{
	"gpt-4o":           "openai",
	"gpt-4o-mini":      "openai",
	"claude-3-sonnet":  "anthropic",
	"claude-3-haiku":   "anthropic",
	"gemini-2.5-flash": "google",
}

var hostProvider = map[string]string{
	"api.openai.com":                    "openai",
	"api.anthropic.com":                 "anthropic",
	"generativelanguage.googleapis.com": "google",
}

func GetAPIKey(host string) (string, error) {
	switch host {
	case "api.openai.com":
		return os.Getenv("OPENAI_API_KEY"), nil
	case "api.anthropic.com":
		return os.Getenv("ANTHROPIC_API_KEY"), nil
	case "generativelanguage.googleapis.com":
		return os.Getenv("GOOGLE_API_KEY"), nil
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

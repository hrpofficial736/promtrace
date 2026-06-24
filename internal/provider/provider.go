package provider

import (
	"fmt"
	"net/http"
	"os"

	"github.com/hrpofficial736/promtrace/internal/logger"
)

type ProviderConfig struct {
	Name       string // "openai", "anthropic", "gemini"
	Host       string // "api.openai.com"
	APIKeyEnv  string // env var name: "OPENAI_API_KEY"
	AuthStyle  string // "bearer", "header", "query"
	AuthHeader string // header name if applicable: "Authorization", "x-api-key"
	ModelIn    string // where the model lives: "body", "url"
}

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
	switch hostToProvider[host] {
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

var providers = map[string]*ProviderConfig{
	"openai": {
		Name:       "openai",
		Host:       "api.openai.com",
		APIKeyEnv:  "OPENAI_API_KEY",
		AuthStyle:  "bearer", // Authorization: Bearer <key>
		AuthHeader: "Authorization",
		ModelIn:    "body",
	},
	"anthropic": {
		Name:       "anthropic",
		Host:       "api.anthropic.com",
		APIKeyEnv:  "ANTHROPIC_API_KEY",
		AuthStyle:  "header", // x-api-key: <key>
		AuthHeader: "x-api-key",
		ModelIn:    "body",
	},
	"gemini": {
		Name:       "gemini",
		Host:       "generativelanguage.googleapis.com",
		APIKeyEnv:  "GEMINI_API_KEY",
		AuthStyle:  "query", // ?key=<key>
		AuthHeader: "",
		ModelIn:    "url", // model is in URL path
	},
}

var modelProvider = map[string]string{
	"gpt-4o":           "openai",
	"gpt-4o-mini":      "openai",
	"claude-3-sonnet":  "anthropic",
	"claude-3-haiku":   "anthropic",
	"gemini-2.5-flash": "gemini",
}

var hostToProvider = map[string]string{
	"api.openai.com":                    "openai",
	"api.anthropic.com":                 "anthropic",
	"generativelanguage.googleapis.com": "gemini",
}

func GetProviderByHost(host string) *ProviderConfig {
	return providers[hostToProvider[host]]
}

func ValidateModel(modelFlag, traceHost string) error {
	if hostToProvider[traceHost] != modelProvider[modelFlag] {
		return fmt.Errorf("cannot replay: trace was sent to %s but model %s belongs to %s", traceHost, modelFlag, modelProvider[modelFlag])
	}
	return nil
}

func SetAuth(req *http.Request, host string) error {
	p := GetProviderByHost(host)

	if p == nil {
		return fmt.Errorf("unknown provider for the host %s", host)
	}

	key := os.Getenv(p.APIKeyEnv)

	if key == "" {
		return fmt.Errorf("env var %s not set", p.APIKeyEnv)
	}

	switch p.AuthStyle {
	case "bearer":
		req.Header.Set(p.AuthHeader, "Bearer "+key)
	case "header":
		req.Header.Set(p.AuthHeader, key)
	case "query":
		q := req.URL.Query()
		q.Set("key", key)
		req.URL.RawQuery = q.Encode()
	}

	return nil
}

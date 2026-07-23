package provider

import (
	"fmt"
	"net/http"
	"os"

	"github.com/hrpofficial736/promtrace/pkg/costable"
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
	Cost         float64
}

type Extractor interface {
	ExtractRequest(body []byte, path string) (model, systemPrompt, userPrompt string)
	ExtractResponse(body []byte) (response string, inputTokens, outputTokens int)
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
	// ── OpenAI: GPT-5.x ──────────────────────────────────────────────────────
	"gpt-5.6-sol":   "openai",
	"gpt-5.6-terra": "openai",
	"gpt-5.6-luna":  "openai",
	"gpt-5.5":       "openai",
	"gpt-5.5-pro":   "openai",
	"gpt-5.4":       "openai",
	"gpt-5.4-mini":  "openai",
	"gpt-5.4-nano":  "openai",
	"gpt-5.4-pro":   "openai",
	"gpt-5.2":       "openai",
	"gpt-5.2-pro":   "openai",
	"gpt-5.1":       "openai",
	"gpt-5":         "openai",
	"gpt-5-mini":    "openai",
	"gpt-5-nano":    "openai",
	"gpt-5-pro":     "openai",

	// ── OpenAI: GPT-4.x ──────────────────────────────────────────────────────
	"gpt-4.1":                   "openai",
	"gpt-4.1-mini":              "openai",
	"gpt-4.1-nano":              "openai",
	"gpt-4o":                    "openai",
	"gpt-4o-2024-05-13":         "openai",
	"gpt-4o-mini":               "openai",
	"gpt-4-turbo-2024-04-09":    "openai",
	"gpt-4-0125-preview":        "openai",
	"gpt-4-1106-preview":        "openai",
	"gpt-4-1106-vision-preview": "openai",
	"gpt-4-0613":                "openai",
	"gpt-4-0314":                "openai",
	"gpt-4-32k":                 "openai",

	// ── OpenAI: o-series (reasoning) ─────────────────────────────────────────
	"o1":                    "openai",
	"o1-pro":                "openai",
	"o1-mini":               "openai",
	"o3":                    "openai",
	"o3-pro":                "openai",
	"o3-mini":               "openai",
	"o4-mini":               "openai",
	"o3-deep-research":      "openai",
	"o4-mini-deep-research": "openai",
	"computer-use-preview":  "openai",

	// ── OpenAI: GPT-3.5 (legacy) ─────────────────────────────────────────────
	"gpt-3.5-turbo":          "openai",
	"gpt-3.5-turbo-0125":     "openai",
	"gpt-3.5-turbo-1106":     "openai",
	"gpt-3.5-turbo-0613":     "openai",
	"gpt-3.5-turbo-0301":     "openai",
	"gpt-3.5-turbo-instruct": "openai",
	"gpt-3.5-turbo-16k-0613": "openai",

	// ── Anthropic: Opus ──────────────────────────────────────────────────────
	"claude-opus-4-8": "anthropic",
	"claude-opus-4-7": "anthropic",
	"claude-opus-4-6": "anthropic",
	"claude-opus-4-5": "anthropic",
	"claude-opus-4-1": "anthropic",
	"claude-opus-4-0": "anthropic",

	// ── Anthropic: Sonnet ────────────────────────────────────────────────────
	"claude-sonnet-5":            "anthropic",
	"claude-sonnet-4-6":          "anthropic",
	"claude-sonnet-4-5":          "anthropic",
	"claude-sonnet-4-0":          "anthropic",
	"claude-3-5-sonnet":          "anthropic",
	"claude-3-5-sonnet-20241022": "anthropic",

	// ── Anthropic: Haiku ─────────────────────────────────────────────────────
	"claude-haiku-4-5":        "anthropic",
	"claude-haiku-3-5":        "anthropic",
	"claude-3-haiku":          "anthropic",
	"claude-3-haiku-20240307": "anthropic",

	// ── Anthropic: other ─────────────────────────────────────────────────────
	"claude-fable-5":  "anthropic",
	"claude-mythos-5": "anthropic",

	// ── Gemini: 3.x ──────────────────────────────────────────────────────────
	"gemini-3.6-flash":                   "gemini",
	"gemini-3.5-flash":                   "gemini",
	"gemini-3.5-flash-lite":              "gemini",
	"gemini-3.1-flash-lite":              "gemini",
	"gemini-3.1-pro-preview":             "gemini",
	"gemini-3.1-pro-preview-customtools": "gemini",
	"gemini-3-flash-preview":             "gemini",
	"gemini-3.1-flash-live-preview":      "gemini",
	"gemini-omni-flash-preview":          "gemini",

	// ── Gemini: 2.5 ──────────────────────────────────────────────────────────
	"gemini-2.5-pro":                                "gemini",
	"gemini-2.5-flash":                              "gemini",
	"gemini-2.5-flash-lite":                         "gemini",
	"gemini-2.5-flash-lite-preview-09-2025":         "gemini",
	"gemini-2.5-flash-image":                        "gemini",
	"gemini-2.5-flash-native-audio-preview-12-2025": "gemini",
	"gemini-2.5-flash-preview-tts":                  "gemini",
	"gemini-2.5-pro-preview-tts":                    "gemini",
	"gemini-2.5-computer-use-preview-10-2025":       "gemini",
	"gemini-robotics-er-1.6-preview":                "gemini",
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
	isModelValid := costable.ValidateModel(modelFlag)
	if !isModelValid {
		return fmt.Errorf("invalid model, please replay the request with a valid model")
	}
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

package provider

import (
	"os"
	"testing"
)

func loadFixture(t *testing.T, path string) []byte {
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to load fixture %s: %v", path, err)
	}
	return data
}

func TestOpenAIExtractRequest(t *testing.T) {
	body := loadFixture(t, "../../testdata/fixtures/openai_request.json")

	ext := &OpenAIExtractor{}

	model, sysPrompt, userPrompt := ext.ExtractRequest(body, "/v1/chat/completions")

	if model != "gpt-4o" {
		t.Errorf("model: got %s, want gpt-4o", model)
	}
	if sysPrompt != "You are a helpful assistant" {
		t.Errorf("system prompt: got %s", sysPrompt)
	}
	if userPrompt != "What is Go?" {
		t.Errorf("user prompt: got %s", userPrompt)
	}
}

func TestOpenAIExtractResponse(t *testing.T) {
	body := loadFixture(t, "../../testdata/fixtures/openai_response.json")
	ext := &OpenAIExtractor{}
	response, inTok, outTok := ext.ExtractResponse(body)
	if response != "Go is a programming language." {
		t.Errorf("response: got %s", response)
	}
	if inTok != 15 {
		t.Errorf("input tokens: got %d, want 15", inTok)
	}
	if outTok != 8 {
		t.Errorf("output tokens: got %d, want 8", outTok)
	}
}

// --- Gemini ---

func TestGeminiExtractRequest(t *testing.T) {
	body := loadFixture(t, "../../testdata/fixtures/gemini_request.json")
	ext := &GeminiExtractor{}

	model, sysPrompt, userPrompt := ext.ExtractRequest(body, "/v1beta/models/gemini-2.5-flash:generateContent")

	if model != "gemini-2.5-flash" {
		t.Errorf("model: got %s, want gemini-2.5-flash", model)
	}
	if sysPrompt != "You are a Go expert" {
		t.Errorf("system prompt: got %s", sysPrompt)
	}
	if userPrompt != "Explain goroutines" {
		t.Errorf("user prompt: got %s", userPrompt)
	}
}

func TestGeminiExtractResponse(t *testing.T) {
	body := loadFixture(t, "../../testdata/fixtures/gemini_response.json")
	ext := &GeminiExtractor{}

	response, inTok, outTok := ext.ExtractResponse(body)

	if response != "Goroutines are lightweight threads managed by the Go runtime." {
		t.Errorf("response: got %s", response)
	}
	if inTok != 12 {
		t.Errorf("input tokens: got %d, want 12", inTok)
	}
	if outTok != 10 {
		t.Errorf("output tokens: got %d, want 10", outTok)
	}
}

// --- Anthropic ---

func TestAnthropicExtractRequest(t *testing.T) {
	body := loadFixture(t, "../../testdata/fixtures/anthropic_request.json")
	ext := &AnthropicExtractor{}

	model, sysPrompt, userPrompt := ext.ExtractRequest(body, "/v1/messages")

	if model != "claude-3-haiku" {
		t.Errorf("model: got %s, want claude-3-haiku", model)
	}
	if sysPrompt != "You are a concise assistant" {
		t.Errorf("system prompt: got %s", sysPrompt)
	}
	if userPrompt != "What is Rust?" {
		t.Errorf("user prompt: got %s", userPrompt)
	}
}

func TestAnthropicExtractResponse(t *testing.T) {
	body := loadFixture(t, "../../testdata/fixtures/anthropic_response.json")
	ext := &AnthropicExtractor{}

	response, inTok, outTok := ext.ExtractResponse(body)

	if response != "Rust is a systems programming language focused on safety." {
		t.Errorf("response: got %s", response)
	}
	if inTok != 18 {
		t.Errorf("input tokens: got %d, want 18", inTok)
	}
	if outTok != 11 {
		t.Errorf("output tokens: got %d, want 11", outTok)
	}
}

func TestValidateModelSameProvider(t *testing.T) {
	err := ValidateModel("gpt-4o", "api.openai.com")
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestValidateModelCrossProvider(t *testing.T) {
	err := ValidateModel("claude-3-sonnet", "api.openai.com")
	if err == nil {
		t.Error("expected error for cross-provider model")
	}
}

func TestGetExtractorUnknownHost(t *testing.T) {
	ext := GetExtractor("unknown.com")
	if ext != nil {
		t.Error("expected nil for unknown host")
	}
}

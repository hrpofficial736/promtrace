package provider

import (
	"encoding/json"
)

type anthropicRequest struct {
	Model    string `json:"model"`
	System   string `json:"system"`
	Messages []struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	} `json:"messages"`
}

type anthropicResponse struct {
	Content []struct {
		Text string `json:"text"`
	} `json:"content"`
	Usage struct {
		InputTokens  int `json:"input_tokens"`
		OutputTokens int `json:"output_tokens"`
	} `json:"usage"`
}

type AnthropicExtractor struct{}

func (e *AnthropicExtractor) ExtractRequest(body []byte, path string) (string, string, string) {
	var req *anthropicRequest

	if err := json.Unmarshal(body, &req); err != nil {
		return "", "", ""
	}
	var userPrompt string
	for _, msg := range req.Messages {
		if msg.Role == "user" {
			userPrompt = msg.Content
		}
	}

	return req.Model, req.System, userPrompt
}

func (e *AnthropicExtractor) ExtractResponse(body []byte) (string, int, int) {
	var res *anthropicResponse

	if err := json.Unmarshal(body, &res); err != nil {
		return "", 0, 0
	}

	if len(res.Content) > 0 {
		return res.Content[0].Text, res.Usage.InputTokens, res.Usage.OutputTokens
	}

	return "", 0, 0
}

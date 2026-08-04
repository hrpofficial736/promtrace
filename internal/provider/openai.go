package provider

import (
	"encoding/json"
)

type openaiRequest struct {
	Model    string `json:"model"`
	Messages []struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	} `json:"messages"`
}

type openaiResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
	} `json:"usage"`
}

type OpenAIExtractor struct{}

func (e *OpenAIExtractor) ExtractRequest(body []byte, path string) (string, string, string) {
	var req *openaiRequest

	if err := json.Unmarshal(body, &req); err != nil {
		return "", "", ""
	}
	var sysPrompt, userPrompt string
	for _, msg := range req.Messages {
		switch msg.Role {
		case "system":
			sysPrompt = msg.Content
		case "user":
			userPrompt = msg.Content
		}
	}

	return req.Model, sysPrompt, userPrompt
}

func (e *OpenAIExtractor) ExtractResponse(body []byte) (string, int, int) {
	var res *openaiResponse

	if err := json.Unmarshal(body, &res); err != nil {
		return "", 0, 0
	}

	if len(res.Choices) > 0 {
		return res.Choices[0].Message.Content, res.Usage.PromptTokens, res.Usage.CompletionTokens
	}

	return "", 0, 0
}

package provider

import (
	"encoding/json"
	"strings"

	"github.com/hrpofficial736/promtrace/internal/logger"
)

type geminiRequest struct {
	SystemInstruction struct {
		Parts []struct {
			Text string `json:"text"`
		} `json:"parts"`
	} `json:"system_instruction"`
	Contents []struct {
		Parts []struct {
			Text string `json:"text"`
		} `json:"parts"`
	} `json:"contents"`
}
type geminiResponse struct {
	Candidates []struct {
		Content struct {
			Parts []struct {
				Text string `json:"text"`
			} `json:"parts"`
		} `json:"content"`
	} `json:"candidates"`
	UsageMetadata struct {
		PromptTokenCount     int `json:"promptTokenCount"`
		CandidatesTokenCount int `json:"candidatesTokenCount"`
	} `json:"usageMetadata"`
}

type GeminiExtractor struct{}

func (e *GeminiExtractor) ExtractRequest(body []byte, path string) (string, string, string) {
	var req *geminiRequest

	if err := json.Unmarshal(body, &req); err != nil {
		logger.Log.Error("error", "error in gemini extractor", err)
		return "", "", ""
	}

	segments := strings.Split(path, "/")

	var modelName string
	for i, segment := range segments {
		if segment == "models" {
			modelName = strings.Split(segments[i+1], ":")[0]
		}
	}

	if len(req.SystemInstruction.Parts) <= 0 {
		return modelName, "", req.Contents[len(req.Contents)-1].Parts[0].Text
	}

	return modelName, req.SystemInstruction.Parts[0].Text, req.Contents[len(req.Contents)-1].Parts[0].Text
}

func (e *GeminiExtractor) ExtractResponse(body []byte) (string, int, int) {
	var res *geminiResponse

	if err := json.Unmarshal(body, &res); err != nil {
		logger.Log.Error("error", "error in gemini response extractor", err)
		return "", 0, 0
	}

	if len(res.Candidates) > 0 {
		return res.Candidates[0].Content.Parts[0].Text, res.UsageMetadata.PromptTokenCount, res.UsageMetadata.CandidatesTokenCount
	}

	return "", 0, 0
}

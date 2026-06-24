package replay

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/hrpofficial736/promtrace/internal/logger"
	"github.com/hrpofficial736/promtrace/internal/provider"
	"github.com/hrpofficial736/promtrace/internal/store"
)

func ReplayRequest(trace *store.Trace, modelFlag string) (*http.Response, error) {

	bodyBytes := []byte(trace.RequestBody)

	if modelFlag != "" {
		var body map[string]any

		json.Unmarshal(bodyBytes, &body)

		p := provider.GetProviderByHost(trace.Host)
		if p.ModelIn == "body" {
			body["model"] = modelFlag
		} else if p.ModelIn == "url" {
			oldPathParts := strings.Split(trace.Path, "/")
			newPathParts := oldPathParts[:len(oldPathParts)-1]

			newPathParts = append(newPathParts, modelFlag+":generateContent")

			trace.Path = strings.Join(newPathParts, "/")
		}

		bodyBytes, _ = json.Marshal(body)
	}

	url := "https://" + trace.Host + trace.Path

	req, _ := http.NewRequest(trace.Method, url, bytes.NewReader(bodyBytes))

	req.Header.Set("Content-Type", "application/json")

	err := provider.SetAuth(req, trace.Host)

	if err != nil {
		logger.Log.Error("error while getting api key", "error", err)
		return nil, err
	}

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	return client.Do(req)
}

package replay

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/hrpofficial736/promtrace/internal/provider"
	"github.com/hrpofficial736/promtrace/internal/store"
)

func ReplayRequest(trace *store.Trace, modelFlag string) (*http.Response, error) {

	bodyBytes := []byte(trace.RequestBody)

	if modelFlag != "" {
		var body map[string]any

		err := json.Unmarshal(bodyBytes, &body)
		if err != nil {
			return nil, err
		}

		p := provider.GetProviderByHost(trace.Host)
		switch p.ModelIn {
		case "body":
			body["model"] = modelFlag
		case "url":
			oldPathParts := strings.Split(trace.Path, "/")
			newPathParts := oldPathParts[:len(oldPathParts)-1]

			newPathParts = append(newPathParts, modelFlag+":generateContent")

			trace.Path = strings.Join(newPathParts, "/")
		}

		bodyBytes, _ = json.Marshal(body)
	}

	url := "https://" + trace.Host + trace.Path

	req, err := http.NewRequest(trace.Method, url, bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	q := req.URL.Query()
	q.Del("key")
	req.URL.RawQuery = q.Encode()
	req.Header.Del("Authorization")
	req.Header.Del("x-api-key")

	err = provider.SetAuth(req, trace.Host)

	if err != nil {
		return nil, err
	}

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	return client.Do(req)
}

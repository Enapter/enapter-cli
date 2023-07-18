package publichttp

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

type TelemetryAPI struct {
	client *Client
}

type NowQuery struct {
	Devices map[string][]string
}

type NowResponse struct {
	Devices map[string]TelemetryByName `json:"devices"`
	Errors  []Error                    `json:"errors"`
}

type TelemetryByName map[string]TelemetrySnapshot

type TelemetrySnapshot struct {
	Value     interface{} `json:"value"`
	Timestamp int64       `json:"timestamp"`
}

func (t *TelemetryAPI) Now(ctx context.Context, query NowQuery) (NowResponse, error) {
	const path = "/telemetry/v1/now"
	url, err := url.Parse(path)
	if err != nil {
		panic(err)
	}

	values := url.Query()
	for deviceID, items := range query.Devices {
		key := fmt.Sprintf("devices[%s]", deviceID)
		value := strings.Join(items, ",")
		values.Set(key, value)
	}
	url.RawQuery = values.Encode()

	req, err := t.client.NewRequestWithContext(ctx, http.MethodGet, url.String(), nil)
	if err != nil {
		return NowResponse{}, fmt.Errorf("create request: %w", err)
	}

	resp, err := t.client.Do(req)
	if err != nil {
		return NowResponse{}, fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	var telemetry NowResponse
	if err := json.NewDecoder(resp.Body).Decode(&telemetry); err != nil {
		return NowResponse{}, fmt.Errorf("unmarshal request: %w", err)
	}

	return telemetry, nil
}

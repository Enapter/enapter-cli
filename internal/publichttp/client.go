package publichttp

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const defaultBaseURL = "https://api.enapter.com"

type Client struct {
	baseURL *url.URL
	client  *http.Client

	// Services used for talking to different parts of the Enapter API.
	Commands  CommandsAPI
	Telemetry TelemetryAPI
}

func NewClient(httpClient *http.Client) *Client {
	c, err := NewClientWithURL(httpClient, defaultBaseURL)
	if err != nil {
		panic(err)
	}
	return c
}

func NewClientWithURL(httpClient *http.Client, baseURL string) (*Client, error) {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	u, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}

	c := &Client{baseURL: u, client: httpClient}
	c.Commands = CommandsAPI{client: c}
	c.Telemetry = TelemetryAPI{client: c}
	return c, nil
}

func (c *Client) NewRequest(method, path string, body io.Reader) (*http.Request, error) {
	return c.NewRequestWithContext(context.Background(), method, path, body)
}

func (c *Client) NewRequestWithContext(
	ctx context.Context, method, path string, body io.Reader,
) (*http.Request, error) {
	u, err := c.baseURL.Parse(path)
	if err != nil {
		return nil, fmt.Errorf("parse url: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, method, u.String(), body)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	return req, err
}

func (c *Client) Do(req *http.Request) (*http.Response, error) {
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= http.StatusMultipleChoices {
		defer resp.Body.Close()

		responseError, err := c.parseResponseError(resp)
		if err != nil {
			return nil, err
		}
		return nil, responseError
	}

	return resp, nil
}

func (c *Client) parseResponseError(r *http.Response) (ResponseError, error) {
	var errors ResponseError

	if r.Body != http.NoBody {
		if err := json.NewDecoder(r.Body).Decode(&errors); err != nil {
			return ResponseError{}, fmt.Errorf("unmarshal body: %w", err)
		}
	}

	errors.StatusCode = r.StatusCode
	if retryAfter := r.Header.Get("Retry-After"); retryAfter != "" {
		duration, err := strconv.Atoi(retryAfter)
		if err != nil {
			return ResponseError{}, fmt.Errorf("parse Retry-After: %w", err)
		}
		errors.RetryAfter = time.Duration(duration) * time.Second
	}

	return errors, nil
}

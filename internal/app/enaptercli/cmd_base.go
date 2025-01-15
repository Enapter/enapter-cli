package enaptercli

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	"github.com/urfave/cli/v2"
)

type cmdBase struct {
	verbose    bool
	token      string
	apiHost    string
	writer     io.Writer
	httpClient *http.Client
}

func (c *cmdBase) Flags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:        "token",
			Usage:       "Enapter API token",
			EnvVars:     []string{"ENAPTER3_API_TOKEN"},
			Hidden:      true,
			Destination: &c.token,
			Category:    "HTTP API Configuration:",
		},
		&cli.StringFlag{
			Name:        "api-host",
			Usage:       "Override API endpoint",
			EnvVars:     []string{"ENAPTER3_API_HOST"},
			Hidden:      true,
			Value:       "https://api.enapter.com",
			Destination: &c.apiHost,
			Category:    "HTTP API Configuration:",
			Action: func(_ *cli.Context, v string) error {
				c.apiHost = strings.TrimSuffix(v, "/")
				return nil
			},
		},
		&cli.BoolFlag{
			Name:        "verbose",
			Usage:       "log extra details about operation",
			Destination: &c.verbose,
		},
	}
}

func (c *cmdBase) Before(cliCtx *cli.Context) error {
	if cliCtx.String("token") == "" {
		return errAPITokenMissed
	}
	c.writer = cliCtx.App.Writer
	c.httpClient = http.DefaultClient
	return nil
}

func (c *cmdBase) HelpTemplate() string {
	return cli.CommandHelpTemplate + `
ENVIRONMENT VARIABLES:
   ENAPTER3_API_TOKEN  Enapter API access token
   ENAPTER3_API_HOST   Enapter API base URL (https://api.enapter.com by default)

`
}

type doHTTPRequestParams struct {
	Method        string
	Path          string
	Query         url.Values
	Body          io.Reader
	ContentType   string
	RespProcessor func(*http.Response) error
}

func (c *cmdBase) doHTTPRequest(ctx context.Context, p doHTTPRequestParams) error {
	req, err := http.NewRequestWithContext(ctx, p.Method, c.apiHost+"/v3"+p.Path, p.Body)
	if err != nil {
		return fmt.Errorf("build http request: %w", err)
	}

	req.Header.Add("X-Enapter-Auth-Token", c.token)
	req.Header.Set("Content-Type", p.ContentType)
	req.URL.RawQuery = p.Query.Encode()

	if c.verbose {
		bb := &bytes.Buffer{}
		if _, err := io.Copy(bb, req.Body); err != nil {
			return fmt.Errorf("reading body for verbose log: %w", err)
		}
		if err := req.Body.Close(); err != nil {
			return fmt.Errorf("closing body for verbose log: %w", err)
		}
		req.Body = io.NopCloser(bb)

		bodyStr := bb.String()
		if p.ContentType != contentTypeJSON {
			bodyStr = base64.RawStdEncoding.EncodeToString(bb.Bytes())
		}

		fmt.Fprintf(c.writer, "== Do http request %s %s\n", p.Method, req.URL.String())
		fmt.Fprintf(c.writer, "=== Begin body\n%s\n=== End body\n", bodyStr)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("do http request: %w", err)
	}
	defer resp.Body.Close()

	if p.RespProcessor == nil {
		return c.defaultRespProcessor(resp)
	}
	return p.RespProcessor(resp)
}

func (c *cmdBase) dialWebSocket(ctx context.Context, path string) (*websocket.Conn, error) {
	url, err := url.Parse(c.apiHost + "/v3" + path)
	if err != nil {
		return nil, fmt.Errorf("parse url: %w", err)
	}

	switch url.Scheme {
	case "https":
		url.Scheme = "wss"
	case "http":
		url.Scheme = "ws"
	}

	headers := make(http.Header)
	headers.Add("X-Enapter-Auth-Token", c.token)

	if c.verbose {
		fmt.Fprintf(c.writer, "== Dialing WebSocket at %s\n", url.String())
	}

	const timeout = 5 * time.Second
	dialer := websocket.Dialer{
		HandshakeTimeout: timeout,
	}

	//nolint:bodyclose // body should be closed by callers
	conn, resp, err := dialer.DialContext(ctx, url.String(), headers)
	if err != nil {
		return nil, fmt.Errorf("dial: %w", err)
	}
	if resp.StatusCode != http.StatusSwitchingProtocols {
		return nil, cli.Exit(parseRespErrorMessage(resp), 1)
	}

	return conn, nil
}

func (c *cmdBase) defaultRespProcessor(resp *http.Response) error {
	if resp.StatusCode != http.StatusOK {
		return cli.Exit(parseRespErrorMessage(resp), 1)
	}

	n, _ := io.Copy(c.writer, resp.Body)
	if n == 0 {
		_, _ = io.WriteString(c.writer, "Request finished without body\n")
	}

	return nil
}

func okRespBodyProcessor(fn func(body io.Reader) error) func(resp *http.Response) error {
	return func(resp *http.Response) error {
		if resp.StatusCode != http.StatusOK {
			return cli.Exit(parseRespErrorMessage(resp), 1)
		}
		return fn(resp.Body)
	}
}

func parseRespErrorMessage(resp *http.Response) string {
	var errs struct {
		Errors []struct {
			Message string `json:"message"`
		} `json:"errors"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&errs); err != nil {
		if !errors.Is(err, io.EOF) {
			return fmt.Sprintf("parse error response: %s", err)
		}
	}

	if len(errs.Errors) > 0 {
		msg := errs.Errors[0].Message
		if len(msg) > 0 {
			return msg
		}
	}

	return fmt.Sprintf("request finished with HTTP status %q, but without error message", resp.Status)
}

package enaptercli

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"slices"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	"github.com/urfave/cli/v2"
)

type cmdBase struct {
	verbose          bool
	token            string
	apiHost          string
	apiAllowInsecure bool
	writer           io.Writer
	httpClient       *http.Client
}

func (c *cmdBase) Flags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:        "token",
			Usage:       "Enapter API token",
			EnvVars:     []string{"ENAPTER3_API_TOKEN"},
			Hidden:      true,
			Destination: &c.token,
		},
		&cli.StringFlag{
			Name:        "api-url",
			Usage:       "override API base URL",
			EnvVars:     []string{"ENAPTER3_API_URL"},
			Value:       "https://api.enapter.com",
			Destination: &c.apiHost,
			Action: func(_ *cli.Context, v string) error {
				c.apiHost = strings.TrimSuffix(v, "/")
				return nil
			},
		},
		&cli.BoolFlag{
			Name:        "api-allow-insecure",
			Usage:       "allow insecure connections to Enapter API",
			EnvVars:     []string{"ENAPTER3_API_ALLOW_INSECURE"},
			Destination: &c.apiAllowInsecure,
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
	c.httpClient = &http.Client{
		Transport: &http.Transport{
			//nolint:gosec // This is needed to allow self-signed certificates on Gateway.
			TLSClientConfig: &tls.Config{InsecureSkipVerify: c.apiAllowInsecure},
		},
	}

	return nil
}

func (c *cmdBase) HelpTemplate() string {
	return cli.CommandHelpTemplate + `
ENVIRONMENT VARIABLES:
   ENAPTER3_API_TOKEN          Enapter API access token
   ENAPTER3_API_URL            Enapter API base URL (https://api.enapter.com by default)
   ENAPTER3_API_ALLOW_INSECURE Allow insecure connections to Enapter API (default false)

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
		bodyStr, err := getRequestBodyString(req, p.ContentType)
		if err != nil {
			return err
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

type runWebSocketParams struct {
	Path          string
	Query         url.Values
	RespProcessor func(io.Reader) error
}

func (c *cmdBase) runWebSocket(ctx context.Context, p runWebSocketParams) error {
	for retry := false; ; retry = true {
		if retry {
			fmt.Fprintln(c.writer, "Reconnecting...")
			time.Sleep(time.Second)
		}

		conn, err := c.dialWebSocket(ctx, p.Path, p.Query)
		if err != nil {
			select {
			case <-ctx.Done():
				return nil
			default:
				fmt.Fprintln(c.writer, "Failed to retrieve data:", err)
				continue
			}
		}
		fmt.Fprintln(c.writer, "Connection established")

		closeCh := make(chan struct{})
		go func() {
			select {
			case <-ctx.Done():
			case <-closeCh:
			}
			conn.Close()
		}()

		if err := c.readWebSocket(conn, p.RespProcessor); err != nil {
			select {
			case <-ctx.Done():
				return nil
			default:
				fmt.Fprintln(c.writer, "Failed to retrieve data:", err)
				close(closeCh)
			}
		}
	}
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

func (c *cmdBase) dialWebSocket(
	ctx context.Context, path string, query url.Values,
) (*websocket.Conn, error) {
	url, err := url.Parse(c.apiHost + "/v3" + path)
	if err != nil {
		return nil, fmt.Errorf("parse url: %w", err)
	}
	url.RawQuery = query.Encode()

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

func (c *cmdBase) readWebSocket(
	conn *websocket.Conn, processor func(io.Reader) error,
) error {
	for {
		_, r, err := conn.NextReader()
		if err != nil {
			return fmt.Errorf("read: %w", err)
		}
		if err := processor(r); err != nil {
			return err
		}
	}
}

func getRequestBodyString(req *http.Request, contentType string) (string, error) {
	if req.Body == nil {
		return "", nil
	}
	bb := &bytes.Buffer{}
	if _, err := io.Copy(bb, req.Body); err != nil {
		return "", fmt.Errorf("reading body for verbose log: %w", err)
	}
	if err := req.Body.Close(); err != nil {
		return "", fmt.Errorf("closing body for verbose log: %w", err)
	}
	req.Body = io.NopCloser(bb)

	if contentType != contentTypeJSON {
		return base64.RawStdEncoding.EncodeToString(bb.Bytes()), nil
	}

	return bb.String(), nil
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

func validateExpandFlag(cliCtx *cli.Context, supportedFields []string) error {
	for _, field := range cliCtx.StringSlice("expand") {
		if err := validateFlag("expand", field, supportedFields); err != nil {
			return err
		}
	}
	return nil
}

func validateFlag(context, value string, allowedValues []string) error {
	slices.Sort(allowedValues)
	if _, ok := slices.BinarySearch(allowedValues, value); !ok {
		return fmt.Errorf("%w: %s is not supported for %s, should be one of %s",
			errUnsupportedFlagValue, value, context, allowedValues)
	}
	return nil
}

package enaptercli

import (
	"bytes"
	"cmp"
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

	"github.com/enapter/enapter-cli/internal/app/configfile"
)

const defaultURL = "https://api.enapter.com"

type cmdBase struct {
	connName         string
	token            string
	apiURL           string
	siteID           string
	apiAllowInsecure bool
	verbose          bool
	colorize         bool
	writer           io.Writer
	errWriter        io.Writer
	userAgent        string
	httpClient       *http.Client
}

func (c *cmdBase) Flags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:        "connection",
			Usage:       "Name of the connection to use",
			Aliases:     []string{"c"},
			Destination: &c.connName,
		},
		&cli.StringFlag{
			Name:        "token",
			Usage:       "Enapter API token",
			EnvVars:     []string{"ENAPTER3_API_TOKEN"},
			Hidden:      true,
			Destination: &c.token,
		},
		&cli.StringFlag{
			Name:        "api-url",
			Usage:       "Override API base URL",
			EnvVars:     []string{"ENAPTER3_API_URL"},
			Value:       defaultURL,
			Hidden:      true,
			Destination: &c.apiURL,
			Action: func(_ *cli.Context, v string) error {
				c.apiURL = strings.TrimSuffix(v, "/")
				return nil
			},
		},
		&cli.BoolFlag{
			Name:        "api-allow-insecure",
			Usage:       "Allow insecure connections to the Enapter API",
			EnvVars:     []string{"ENAPTER3_API_ALLOW_INSECURE"},
			Destination: &c.apiAllowInsecure,
		},
		&cli.BoolFlag{
			Name:        "verbose",
			Usage:       "Log extra details about the operation",
			Destination: &c.verbose,
		},
	}
}

func (c *cmdBase) Before(cliCtx *cli.Context) error {
	if err := c.setupCredentials(cliCtx); err != nil {
		return err
	}

	c.writer = cliCtx.App.Writer
	c.errWriter = cliCtx.App.ErrWriter
	c.colorize = colorsSupported(c.writer)

	c.userAgent = "enapter-cli/" + cliCtx.App.Version
	c.httpClient = &http.Client{
		Transport: &http.Transport{
			//nolint:gosec // This is needed to allow self-signed certificates on Gateway.
			TLSClientConfig: &tls.Config{InsecureSkipVerify: c.apiAllowInsecure},
		},
	}

	return nil
}

func (c *cmdBase) setupCredentials(cliCtx *cli.Context) error {
	config, err := configfile.Load()
	if err != nil {
		return err
	}

	if c.connName != "" {
		conn, ok := config.Connections[c.connName]
		if !ok {
			return cli.Exit("Unknown connection name.", 1)
		}
		if cliCtx.IsSet("token") || cliCtx.IsSet("api-url") || cliCtx.IsSet("api-allow-insecure") {
			fmt.Fprintln(cliCtx.App.ErrWriter,
				"WARNING: credentials set via environment variables or flags are ignored.")
		}
		c.token = conn.Token.Value
		c.apiURL = conn.URL
		c.siteID = conn.SiteID
		c.apiAllowInsecure = conn.AllowInsecure
		return nil
	}

	if c.token != "" {
		return nil
	}

	if config.DefaultConn != "" {
		conn, ok := config.Connections[config.DefaultConn]
		if !ok {
			return cli.Exit("Default connection is invalid.", 1)
		}
		c.token = conn.Token.Value
		c.apiURL = conn.URL
		c.siteID = conn.SiteID
		c.apiAllowInsecure = conn.AllowInsecure
		return nil
	}

	return cli.Exit("No connection configured.\n\n"+
		"Please, specify connection using --connection flag.\n\n"+
		"To list available connections:\n$ enapter3 connection list\n\n"+
		"To add a new connection:\n$ enapter3 connection add\n", 1)
}

const enapterAPIEnvVarsHelp = `
ENVIRONMENT VARIABLES:
   ENAPTER3_API_TOKEN          Enapter API access token
   ENAPTER3_API_URL            Enapter API base URL (default: https://api.enapter.com)
   ENAPTER3_API_ALLOW_INSECURE Allow insecure connections to the Enapter API (default: false)

`

func (c *cmdBase) CommandHelpTemplate() string {
	return cli.CommandHelpTemplate + enapterAPIEnvVarsHelp
}

func (c *cmdBase) SubcommandHelpTemplate() string {
	return cli.SubcommandHelpTemplate + enapterAPIEnvVarsHelp
}

func (c *cmdBase) chooseSiteID(cmdSiteID string) (string, error) {
	if cmdSiteID != "" && c.siteID != "" && c.siteID != cmdSiteID {
		return "", errSiteIDMismatch
	}
	siteID := cmp.Or(cmdSiteID, c.siteID)
	if siteID == "" {
		return "", errSiteIDMissing
	}
	return siteID, nil
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
	req, err := http.NewRequestWithContext(ctx, p.Method, c.apiURL+"/v3"+p.Path, p.Body)
	if err != nil {
		return fmt.Errorf("build http request: %w", err)
	}

	req.Header.Set("X-Enapter-Auth-Token", c.token)
	req.Header.Set("User-Agent", c.userAgent)
	req.Header.Set("Content-Type", p.ContentType)
	req.URL.RawQuery = p.Query.Encode()

	if c.verbose {
		bodyStr, err := getRequestBodyString(req, p.ContentType)
		if err != nil {
			return err
		}

		fmt.Fprintf(c.errWriter, "== Do http request %s %s\n", p.Method, req.URL.String())
		fmt.Fprintf(c.errWriter, "=== Begin body\n%s\n=== End body\n", bodyStr)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		if e := (&tls.CertificateVerificationError{}); errors.As(err, &e) {
			return fmt.Errorf("do http request: %w (try to use --api-allow-insecure)", err)
		}
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
	url, err := url.Parse(c.apiURL + "/v3" + p.Path)
	if err != nil {
		return fmt.Errorf("parse url: %w", err)
	}
	url.RawQuery = p.Query.Encode()

	headers := make(http.Header)
	headers.Set("X-Enapter-Auth-Token", c.token)
	headers.Set("User-Agent", c.userAgent)

	for retry := false; ; retry = true {
		if retry {
			fmt.Fprintln(c.errWriter, "Reconnecting...")
			time.Sleep(time.Second)
		}

		conn, err := c.dialWebSocket(ctx, url, headers)
		if err != nil {
			if e := cli.ExitCoder(nil); errors.As(err, &e) {
				return err
			}
			select {
			case <-ctx.Done():
				return nil
			default:
				fmt.Fprintln(c.errWriter, "Failed to retrieve data:", err)
				continue
			}
		}
		fmt.Fprintln(c.errWriter, "Connection established")

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
				fmt.Fprintln(c.errWriter, "Failed to retrieve data:", err)
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
	ctx context.Context, url *url.URL, headers http.Header,
) (*websocket.Conn, error) {
	const timeout = 5 * time.Second
	dialer := websocket.Dialer{
		HandshakeTimeout: timeout,
		//nolint:gosec // This is needed to allow self-signed certificates on Gateway.
		TLSClientConfig: &tls.Config{InsecureSkipVerify: c.apiAllowInsecure},
	}

	const maxRetries = 2
	for i := 0; i < maxRetries; i++ {
		url.Scheme = websocketScheme(url.Scheme)

		if c.verbose {
			fmt.Fprintf(c.errWriter, "== Dialing WebSocket at %s\n", url.String())
		}

		//nolint:bodyclose // body should be closed by callers
		conn, resp, err := dialer.DialContext(ctx, url.String(), headers)
		if err != nil {
			if loc, err := redirectLocation(resp); err != nil {
				return nil, err
			} else if loc != nil {
				url = loc
				continue
			}
			if e := (&tls.CertificateVerificationError{}); errors.As(err, &e) {
				message := fmt.Sprintf("dial: %v (try to use --api-allow-insecure)", err)
				return nil, cli.Exit(message, 1)
			}
			if resp != nil {
				message := parseRespErrorMessage(resp)
				return nil, fmt.Errorf("dial: %w: %s", err, message)
			}
			return nil, fmt.Errorf("dial: %w", err)
		}

		return conn, nil
	}

	return nil, cli.Exit("Too many redirects", 1)
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
	bodyBytes, _ := io.ReadAll(resp.Body)
	if err := json.Unmarshal(bodyBytes, &errs); err != nil {
		if !errors.Is(err, io.EOF) {
			return fmt.Sprintf("Request finished with HTTP status %q, but body is not valid JSON error response. "+
				"Please, check API URL is correct.\n\nReceived body:\n%s\n", resp.Status, bodyBytes)
		}
	}

	if len(errs.Errors) > 0 {
		msg := errs.Errors[0].Message
		if len(msg) > 0 {
			return msg
		}
	}

	return fmt.Sprintf("Request finished with HTTP status %q, but without error message", resp.Status)
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

func websocketScheme(s string) string {
	switch s {
	case "https":
		return "wss"
	case "http":
		return "ws"
	default:
		return s
	}
}

func redirectLocation(resp *http.Response) (*url.URL, error) {
	if resp == nil {
		return nil, nil
	}
	if resp.StatusCode != http.StatusPermanentRedirect {
		return nil, nil
	}
	location := resp.Header.Get("Location")
	url, err := url.Parse(location)
	if err != nil {
		return nil, fmt.Errorf("parse location: %w", err)
	}
	return url, nil
}

package enaptercli

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/urfave/cli/v2"
)

type cmdBase struct {
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
			EnvVars:     []string{"ENAPTER_API_TOKEN"},
			Hidden:      true,
			Destination: &c.token,
			Category:    "HTTP API Configuration:",
		},
		&cli.StringFlag{
			Name:        "api-host",
			Usage:       "Override API endpoint",
			EnvVars:     []string{"ENAPTER_API_HOST"},
			Hidden:      true,
			Value:       "https://api.enapter.com",
			Destination: &c.apiHost,
			Category:    "HTTP API Configuration:",
			Action: func(_ *cli.Context, v string) error {
				c.apiHost = strings.TrimSuffix(v, "/")
				return nil
			},
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
   ENAPTER_API_TOKEN  Enapter API access token
   ENAPTER_API_HOST   Enapter API base URL (https://api.enapter.com by default)

`
}

type doHTTPRequestParams struct {
	Method        string
	Path          string
	Query         url.Values
	Body          io.Reader
	RespProcessor func(*http.Response) error
}

func (c *cmdBase) doHTTPRequest(ctx context.Context, p doHTTPRequestParams) error {
	req, err := http.NewRequestWithContext(ctx, p.Method, c.apiHost+"/v2"+p.Path, p.Body)
	if err != nil {
		return fmt.Errorf("build http request: %w", err)
	}

	req.Header.Add("X-Enapter-Auth-Token", c.token)
	req.URL.RawQuery = p.Query.Encode()

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

package enaptercli

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/urfave/cli/v3"

	"github.com/enapter/enapter-cli/internal/app/configfile"
)

type cmdConnectionAdd struct {
	name          string
	url           string
	token         string
	siteID        string
	gateway       bool
	allowInsecure bool
}

func buildCmdConnectionAdd() *cli.Command {
	cmd := &cmdConnectionAdd{}
	return &cli.Command{
		Name:  "add",
		Usage: "Add a new connection",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "name",
				Usage:       "Connection name",
				Destination: &cmd.name,
				Required:    true,
			},
			&cli.BoolFlag{
				Name:        "gateway",
				Usage:       "Indicates that the connection is to a Gateway",
				Destination: &cmd.gateway,
			},
			&cli.StringFlag{
				Name:        "url",
				Usage:       "Enapter API base URL",
				Destination: &cmd.url,
				Value:       defaultURL,
			},
			&cli.StringFlag{
				Name:        "token",
				Usage:       "Enapter API access token",
				Destination: &cmd.token,
				Required:    true,
			},
			&cli.StringFlag{
				Name: "site-id",
				Usage: "If specified, the connection will be limited to this site " +
					"(available only for Cloud connections)",
				Destination: &cmd.siteID,
			},
			&cli.BoolFlag{
				Name:        "allow-insecure",
				Usage:       "Allow insecure connections to the Enapter API",
				Destination: &cmd.allowInsecure,
			},
		},
		Action: cmd.do,
	}
}

func (c *cmdConnectionAdd) do(ctx context.Context, cliCmd *cli.Command) error {
	config, err := configfile.Load()
	if err != nil {
		return err
	}

	if _, exists := config.Connections[c.name]; exists {
		return cli.Exit("Connection with the given name already exists.", 1)
	}

	if u, err := url.Parse(c.url); err != nil {
		return cli.Exit("Invalid URL format: "+err.Error()+".", 1)
	} else if u.Scheme != "https" && u.Scheme != "http" {
		return cli.Exit("URL scheme must be http or https.", 1)
	}

	if c.gateway {
		if c.url == defaultURL {
			return cli.Exit("Gateway connections require a custom URL.", 1)
		}
		if c.siteID != "" {
			return cli.Exit("The site-id option cannot be used with gateway connections.", 1)
		}
		siteID, err := c.resolveGatewaySiteID(ctx, cliCmd)
		if err != nil {
			return err
		}
		c.siteID = siteID
	}

	if config.Connections == nil {
		config.Connections = make(map[string]configfile.Connection)
	}
	config.Connections[c.name] = configfile.Connection{
		Gateway:       c.gateway,
		URL:           c.url,
		SiteID:        c.siteID,
		Token:         configfile.Token{Value: c.token},
		AllowInsecure: c.allowInsecure,
	}

	return configfile.Save(config)
}

func (c *cmdConnectionAdd) resolveGatewaySiteID(ctx context.Context, cliCmd *cli.Command) (string, error) {
	client := &http.Client{
		Transport: &http.Transport{
			//nolint:gosec // This is needed to allow self-signed certificates on Gateway.
			TLSClientConfig: &tls.Config{InsecureSkipVerify: c.allowInsecure},
		},
	}

	req, err := http.NewRequestWithContext(
		ctx, http.MethodGet, c.url+"/v3/site", nil)
	if err != nil {
		return "", fmt.Errorf("new http request: %w", err)
	}

	req.Header.Set("X-Enapter-Auth-Token", c.token)
	req.Header.Set("User-Agent", "enapter-cli/"+cliCmd.Root().Version)

	resp, err := client.Do(req)
	if err != nil {
		if e := (&tls.CertificateVerificationError{}); errors.As(err, &e) {
			return "", fmt.Errorf("resolve site: %w (try to use --allow-insecure)", err)
		}
		return "", fmt.Errorf("send http request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", cli.Exit("Unexpected response from Gateway: "+resp.Status+". ", 1)
	}

	var siteResp struct {
		Site struct {
			ID string `json:"id"`
		} `json:"site"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&siteResp); err != nil {
		return "", fmt.Errorf("decode http response: %w", err)
	}

	return siteResp.Site.ID, nil
}

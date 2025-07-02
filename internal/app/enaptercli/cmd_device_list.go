package enaptercli

import (
	"context"
	"net/http"
	"net/url"
	"strings"

	"github.com/urfave/cli/v2"
)

type cmdDevicesList struct {
	cmdDevices
	siteID string
	expand []string
}

func buildCmdDevicesList() *cli.Command {
	cmd := &cmdDevicesList{}
	return &cli.Command{
		Name:               "list",
		Usage:              "List user devices",
		CustomHelpTemplate: cmd.CommandHelpTemplate(),
		Flags:              cmd.Flags(),
		Before:             cmd.Before,
		Action: func(cliCtx *cli.Context) error {
			return cmd.do(cliCtx.Context)
		},
	}
}

func (c *cmdDevicesList) Flags() []cli.Flag {
	flags := c.cmdDevices.Flags()
	return append(flags, &cli.MultiStringFlag{
		Target: &cli.StringSliceFlag{
			Name: "expand",
			Usage: "coma-separated list of expanded device information (supported values: " +
				"connectivity, manifest, properties, communication_info, site)",
		},
		Destination: &c.expand,
	}, &cli.StringFlag{
		Name:        "site-id",
		Usage:       "list devices from this site",
		Destination: &c.siteID,
	})
}

func (c *cmdDevicesList) Before(cliCtx *cli.Context) error {
	if err := c.cmdDevices.Before(cliCtx); err != nil {
		return err
	}
	return c.validateExpandFlag(cliCtx)
}

func (c *cmdDevicesList) do(ctx context.Context) error {
	query := url.Values{}
	if len(c.expand) != 0 {
		query.Set("expand", strings.Join(c.expand, ","))
	}

	if c.siteID != "" {
		query.Set("site_id", c.siteID)
		return c.cmdBase.doHTTPRequest(ctx, doHTTPRequestParams{
			Method: http.MethodGet,
			Path:   "/sites/" + c.siteID + "/devices",
			Query:  query,
		})
	}

	return c.doHTTPRequest(ctx, doHTTPRequestParams{
		Method: http.MethodGet,
		Path:   "",
		Query:  query,
	})
}

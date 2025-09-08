package enaptercli

import (
	"context"
	"net/http"
	"net/url"
	"strings"

	"github.com/urfave/cli/v2"
)

type cmdDeviceList struct {
	cmdDevice
	siteID string
	expand []string
	limit  int
}

func buildCmdDeviceList() *cli.Command {
	cmd := &cmdDeviceList{}
	return &cli.Command{
		Name:               "list",
		Usage:              "List user devices ordered by device ID",
		CustomHelpTemplate: cmd.CommandHelpTemplate(),
		Flags:              cmd.Flags(),
		Before:             cmd.Before,
		Action: func(cliCtx *cli.Context) error {
			return cmd.do(cliCtx.Context)
		},
	}
}

func (c *cmdDeviceList) Flags() []cli.Flag {
	flags := c.cmdDevice.Flags()
	return append(flags, &cli.MultiStringFlag{
		Target: &cli.StringSliceFlag{
			Name: "expand",
			Usage: "coma-separated list of expanded device information (supported values: " +
				strings.Join(c.supportedExpandFields(), ", ") + ")",
		},
		Destination: &c.expand,
	}, &cli.StringFlag{
		Name:        "site-id",
		Usage:       "list devices from this site",
		Destination: &c.siteID,
	}, &cli.IntFlag{
		Name:        "limit",
		Usage:       "maximum number of devices to retrieve",
		Destination: &c.limit,
		DefaultText: "retrieves all",
	})
}

func (c *cmdDeviceList) Before(cliCtx *cli.Context) error {
	if err := c.cmdDevice.Before(cliCtx); err != nil {
		return err
	}
	return c.validateExpandFlag(cliCtx)
}

func (c *cmdDeviceList) do(ctx context.Context) error {
	query := url.Values{}
	if len(c.expand) != 0 {
		query.Set("expand", strings.Join(c.expand, ","))
	}

	doPaginateRequestParams := paginateHTTPRequestParams{
		ObjectName: "devices",
		Limit:      c.limit,
		DoFn:       c.doHTTPRequest,
		BaseParams: doHTTPRequestParams{
			Method: http.MethodGet,
			Path:   "",
			Query:  query,
		},
	}

	if c.siteID != "" {
		doPaginateRequestParams.BaseParams.Query.Set("site_id", c.siteID)
		doPaginateRequestParams.BaseParams.Path = "/sites/" + c.siteID + "/devices"
		doPaginateRequestParams.DoFn = c.cmdBase.doHTTPRequest
	}

	return c.doPaginateRequest(ctx, doPaginateRequestParams)
}

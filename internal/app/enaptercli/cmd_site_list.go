package enaptercli

import (
	"context"
	"net/http"

	"github.com/urfave/cli/v2"
)

type cmdSitesList struct {
	cmdSite
	mySites bool
	limit   int
}

func buildCmdSitesList() *cli.Command {
	cmd := &cmdSitesList{}
	return &cli.Command{
		Name:               "list",
		Usage:              "List user sites",
		CustomHelpTemplate: cmd.CommandHelpTemplate(),
		Flags:              cmd.Flags(),
		Before:             cmd.Before,
		Action: func(cliCtx *cli.Context) error {
			return cmd.do(cliCtx.Context)
		},
	}
}

func (c *cmdSitesList) Flags() []cli.Flag {
	flags := c.cmdSite.Flags()
	return append(flags, &cli.BoolFlag{
		Name:        "my-sites",
		Usage:       "returns only sites where user is owner or installer",
		Destination: &c.mySites,
	}, &cli.IntFlag{
		Name:        "limit",
		Usage:       "maximum number of sites to retrieve",
		Destination: &c.limit,
		DefaultText: "retrieves all",
	})
}

func (c *cmdSitesList) do(ctx context.Context) error {
	doPaginateRequestParams := paginateHTTPRequestParams{
		ObjectName: "sites",
		Limit:      c.limit,
		DoFn:       c.doHTTPRequest,
		BaseParams: doHTTPRequestParams{
			Method: http.MethodGet,
			Path:   "",
		},
	}

	if c.mySites {
		doPaginateRequestParams.BaseParams.Path = "/users/me/sites"
	}

	return c.doPaginateRequest(ctx, doPaginateRequestParams)
}

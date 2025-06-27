package enaptercli

import (
	"context"
	"net/http"

	"github.com/urfave/cli/v2"
)

type cmdSiteGet struct {
	cmdSite
	siteID string
}

func buildCmdSiteGet() *cli.Command {
	cmd := &cmdSiteGet{}
	return &cli.Command{
		Name:               "get",
		Usage:              "Get a site",
		CustomHelpTemplate: cmd.HelpTemplate(),
		Flags:              cmd.Flags(),
		Before:             cmd.Before,
		Action: func(cliCtx *cli.Context) error {
			return cmd.do(cliCtx.Context)
		},
	}
}

func (c *cmdSiteGet) Flags() []cli.Flag {
	flags := c.cmdSite.Flags()
	return append(flags, &cli.StringFlag{
		Name:        "site-id",
		Usage:       "site ID",
		Destination: &c.siteID,
		Required:    true,
	})
}

func (c *cmdSiteGet) do(ctx context.Context) error {
	return c.doHTTPRequest(ctx, doHTTPRequestParams{
		Method: http.MethodGet,
		Path:   "/" + c.siteID,
	})
}

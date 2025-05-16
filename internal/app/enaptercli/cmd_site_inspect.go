package enaptercli

import (
	"context"
	"net/http"

	"github.com/urfave/cli/v2"
)

type cmdSiteInspect struct {
	cmdSite
	siteID string
}

func buildCmdSiteInspect() *cli.Command {
	cmd := &cmdSiteInspect{}
	return &cli.Command{
		Name:               "inspect",
		Usage:              "Inspect a site",
		CustomHelpTemplate: cmd.HelpTemplate(),
		Flags:              cmd.Flags(),
		Before:             cmd.Before,
		Action: func(cliCtx *cli.Context) error {
			return cmd.do(cliCtx.Context)
		},
	}
}

func (c *cmdSiteInspect) Flags() []cli.Flag {
	flags := c.cmdSite.Flags()
	return append(flags, &cli.StringFlag{
		Name:        "site-id",
		Usage:       "site ID",
		Destination: &c.siteID,
		Required:    true,
	})
}

func (c *cmdSiteInspect) do(ctx context.Context) error {
	return c.doHTTPRequest(ctx, doHTTPRequestParams{
		Method: http.MethodGet,
		Path:   "/" + c.siteID,
	})
}

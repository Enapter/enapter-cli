package enaptercli

import (
	"context"
	"net/http"

	"github.com/urfave/cli/v2"
)

type cmdSitesList struct {
	cmdSite
	mySites bool
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
	})
}

func (c *cmdSitesList) do(ctx context.Context) error {
	if c.mySites {
		return c.cmdBase.doHTTPRequest(ctx, doHTTPRequestParams{
			Method: http.MethodGet,
			Path:   "/users/me/sites",
		})
	}

	return c.cmdSite.doHTTPRequest(ctx, doHTTPRequestParams{Method: http.MethodGet})
}

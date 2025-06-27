package enaptercli

import (
	"context"
	"net/http"
	"net/url"
	"strconv"

	"github.com/urfave/cli/v2"
)

type cmdSitesList struct {
	cmdSite
	offset int
	limit  int
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

func (c *cmdSitesList) do(ctx context.Context) error {
	query := url.Values{}
	if c.offset != 0 {
		query.Set("offset", strconv.Itoa(c.offset))
	}
	if c.limit != 0 {
		query.Set("limit", strconv.Itoa(c.limit))
	}

	return c.cmdBase.doHTTPRequest(ctx, doHTTPRequestParams{
		Method: http.MethodGet,
		Path:   "/users/me/sites",
		Query:  query,
	})
}

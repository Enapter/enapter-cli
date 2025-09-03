package enaptercli

import (
	"context"
	"fmt"
	"net/url"

	"github.com/urfave/cli/v2"
)

type cmdSite struct {
	cmdBase
}

func buildCmdSites() *cli.Command {
	cmd := &cmdSite{}
	return &cli.Command{
		Name:               "site",
		Usage:              "Manage sites",
		CustomHelpTemplate: cmd.SubcommandHelpTemplate(),
		Subcommands: []*cli.Command{
			buildCmdSitesList(),
			buildCmdSiteGet(),
		},
	}
}

func (c *cmdSite) doHTTPRequest(ctx context.Context, p doHTTPRequestParams) error {
	path, err := url.JoinPath("/sites", p.Path)
	if err != nil {
		return fmt.Errorf("join path: %w", err)
	}
	p.Path = path
	return c.cmdBase.doHTTPRequest(ctx, p)
}

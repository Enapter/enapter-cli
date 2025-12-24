package enaptercli

import (
	"cmp"
	"context"
	"fmt"
	"net/url"

	"github.com/urfave/cli/v2"
)

type cmdRuleEngine struct {
	cmdBase
	siteID string
}

func buildCmdRuleEngine() *cli.Command {
	cmd := &cmdRuleEngine{}
	return &cli.Command{
		Name:               "rule-engine",
		Usage:              "Manage the rule engine",
		CustomHelpTemplate: cmd.SubcommandHelpTemplate(),
		Subcommands: []*cli.Command{
			buildCmdRuleEngineGet(),
			buildCmdRuleEngineSuspend(),
			buildCmdRuleEngineResume(),
			buildCmdRuleEngineRule(),
		},
	}
}

func (c *cmdRuleEngine) Flags() []cli.Flag {
	flags := c.cmdBase.Flags()
	return append(flags, &cli.StringFlag{
		Name:        "site-id",
		Usage:       "site ID",
		Destination: &c.siteID,
	})
}

func (c *cmdRuleEngine) doHTTPRequest(ctx context.Context, p doHTTPRequestParams) error {
	if c.siteID != "" && c.cmdBase.siteID != "" && c.cmdBase.siteID != c.siteID {
		return errSiteIDMismatch
	}

	siteID := cmp.Or(c.siteID, c.cmdBase.siteID)
	if siteID == "" {
		return errSiteIDMissing
	}

	path, err := url.JoinPath("/sites/", siteID, "/rule_engine", p.Path)
	if err != nil {
		return fmt.Errorf("join path: %w", err)
	}

	p.Path = path
	return c.cmdBase.doHTTPRequest(ctx, p)
}

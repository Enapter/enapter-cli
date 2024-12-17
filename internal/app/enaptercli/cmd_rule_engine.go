package enaptercli

import (
	"context"
	"fmt"
	"net/url"

	"github.com/urfave/cli/v2"
)

type cmdRuleEngine struct {
	cmdBase
}

func buildCmdRuleEngine() *cli.Command {
	return &cli.Command{
		Name:  "rule-engine",
		Usage: "Manage the rule engine",
		Subcommands: []*cli.Command{
			buildCmdRuleEngineInspect(),
			buildCmdRuleEngineSuspend(),
			buildCmdRuleEngineResume(),
			buildCmdRuleEngineRule(),
		},
	}
}

func (c *cmdRuleEngine) doHTTPRequest(ctx context.Context, p doHTTPRequestParams) error {
	path, err := url.JoinPath("/site/rule_engine", p.Path)
	if err != nil {
		return fmt.Errorf("join path: %w", err)
	}
	p.Path = path
	return c.cmdBase.doHTTPRequest(ctx, p)
}

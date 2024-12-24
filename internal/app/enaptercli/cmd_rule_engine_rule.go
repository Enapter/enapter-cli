package enaptercli

import (
	"context"
	"fmt"
	"net/url"

	"github.com/urfave/cli/v2"
)

const (
	ruleRuntimeVersion1 = 1
	ruleRuntimeVersion3 = 3
)

type cmdRuleEngineRule struct {
	cmdRuleEngine
}

func buildCmdRuleEngineRule() *cli.Command {
	return &cli.Command{
		Name:  "rule",
		Usage: "Manage rules",
		Subcommands: []*cli.Command{
			buildCmdRuleEngineRuleCreate(),
			buildCmdRuleEngineRuleDelete(),
			buildCmdRuleEngineRuleDisable(),
			buildCmdRuleEngineRuleEnable(),
			buildCmdRuleEngineRuleInspect(),
			buildCmdRuleEngineRuleList(),
			buildCmdRuleEngineRuleUpdate(),
			buildCmdRuleEngineRuleUpdateScript(),
		},
	}
}

func (c *cmdRuleEngineRule) doHTTPRequest(ctx context.Context, p doHTTPRequestParams) error {
	path, err := url.JoinPath("/rules", p.Path)
	if err != nil {
		return fmt.Errorf("join path: %w", err)
	}
	p.Path = path
	return c.cmdRuleEngine.doHTTPRequest(ctx, p)
}

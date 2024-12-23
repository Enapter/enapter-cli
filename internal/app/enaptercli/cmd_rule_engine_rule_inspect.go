package enaptercli

import (
	"context"
	"net/http"

	"github.com/urfave/cli/v3"
)

type cmdRuleEngineRuleInspect struct {
	cmdRuleEngineRule
}

func buildCmdRuleEngineRuleInspect() *cli.Command {
	cmd := &cmdRuleEngineRuleInspect{}
	return &cli.Command{
		Name:               "inspect",
		Usage:              "Inspect a rule",
		ArgsUsage:          "RULE",
		CustomHelpTemplate: cmd.HelpTemplate(),
		Flags:              cmd.Flags(),
		Before:             cmd.Before,
		Action: func(ctx context.Context, cm *cli.Command) error {
			return cmd.do(ctx, cm.Args().First())
		},
	}
}

func (c *cmdRuleEngineRuleInspect) Before(
	ctx context.Context, cm *cli.Command,
) (context.Context, error) {
	ctx, err := c.cmdRuleEngineRule.Before(ctx, cm)
	if err != nil {
		return nil, err
	}

	if cm.Args().Get(0) == "" {
		return nil, errRequiresAtLeastOneArgument
	}

	return ctx, nil
}

func (c *cmdRuleEngineRuleInspect) do(ctx context.Context, rule string) error {
	return c.doHTTPRequest(ctx, doHTTPRequestParams{
		Method: http.MethodGet,
		Path:   "/" + rule,
	})
}

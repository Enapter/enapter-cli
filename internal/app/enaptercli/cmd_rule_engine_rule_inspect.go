package enaptercli

import (
	"context"
	"net/http"

	"github.com/urfave/cli/v2"
)

type cmdRuleEngineRuleInspect struct {
	cmdRuleEngineRule
}

func buildCmdRuleEngineRuleInspect() *cli.Command {
	cmd := &cmdRuleEngineRuleInspect{}
	return &cli.Command{
		Name:               "inspect",
		Usage:              "Inspect a rule",
		Args:               true,
		ArgsUsage:          "RULE",
		CustomHelpTemplate: cmd.HelpTemplate(),
		Flags:              cmd.Flags(),
		Before:             cmd.Before,
		Action: func(cliCtx *cli.Context) error {
			return cmd.do(cliCtx.Context, cliCtx.Args().First())
		},
	}
}

func (c *cmdRuleEngineRuleInspect) Before(cliCtx *cli.Context) error {
	if err := c.cmdRuleEngineRule.Before(cliCtx); err != nil {
		return err
	}

	if cliCtx.Args().Get(0) == "" {
		return errRequiresAtLeastOneArgument
	}

	return nil
}

func (c *cmdRuleEngineRuleInspect) do(ctx context.Context, rule string) error {
	return c.doHTTPRequest(ctx, doHTTPRequestParams{
		Method: http.MethodGet,
		Path:   "/" + rule,
	})
}

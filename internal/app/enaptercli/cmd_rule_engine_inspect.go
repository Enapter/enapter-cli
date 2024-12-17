package enaptercli

import (
	"context"
	"net/http"

	"github.com/urfave/cli/v2"
)

type cmdRuleEngineInspect struct {
	cmdRuleEngine
}

func buildCmdRuleEngineInspect() *cli.Command {
	cmd := &cmdRuleEngineInspect{}
	return &cli.Command{
		Name:               "inspect",
		Usage:              "Inspect the rule engine",
		CustomHelpTemplate: cmd.HelpTemplate(),
		Flags:              cmd.Flags(),
		Before:             cmd.Before,
		Action: func(cliCtx *cli.Context) error {
			return cmd.do(cliCtx.Context)
		},
	}
}

func (c *cmdRuleEngineInspect) do(ctx context.Context) error {
	return c.doHTTPRequest(ctx, doHTTPRequestParams{
		Method: http.MethodGet,
		Path:   "",
	})
}

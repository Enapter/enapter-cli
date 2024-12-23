package enaptercli

import (
	"context"
	"net/http"

	"github.com/urfave/cli/v3"
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
		Action: func(ctx context.Context, _ *cli.Command) error {
			return cmd.do(ctx)
		},
	}
}

func (c *cmdRuleEngineInspect) do(ctx context.Context) error {
	return c.doHTTPRequest(ctx, doHTTPRequestParams{
		Method: http.MethodGet,
		Path:   "",
	})
}

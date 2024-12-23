package enaptercli

import (
	"context"
	"net/http"

	"github.com/urfave/cli/v3"
)

type cmdRuleEngineRuleList struct {
	cmdRuleEngineRule
}

func buildCmdRuleEngineRuleList() *cli.Command {
	cmd := &cmdRuleEngineRuleList{}
	return &cli.Command{
		Name:               "list",
		Usage:              "List rules",
		CustomHelpTemplate: cmd.HelpTemplate(),
		Flags:              cmd.Flags(),
		Before:             cmd.Before,
		Action: func(ctx context.Context, _ *cli.Command) error {
			return cmd.do(ctx)
		},
	}
}

func (c *cmdRuleEngineRuleList) do(ctx context.Context) error {
	return c.doHTTPRequest(ctx, doHTTPRequestParams{
		Method: http.MethodGet,
		Path:   "",
	})
}

package enaptercli

import (
	"context"
	"net/http"

	"github.com/urfave/cli/v2"
)

type cmdRuleEngineRuleList struct {
	cmdRuleEngineRule
}

func buildCmdRuleEngineRuleList() *cli.Command {
	cmd := &cmdRuleEngineRuleList{}
	return &cli.Command{
		Name:               "list",
		Usage:              "List rules",
		CustomHelpTemplate: cmd.CommandHelpTemplate(),
		Flags:              cmd.Flags(),
		Before:             cmd.Before,
		Action: func(cliCtx *cli.Context) error {
			return cmd.do(cliCtx.Context)
		},
	}
}

func (c *cmdRuleEngineRuleList) do(ctx context.Context) error {
	return c.doHTTPRequest(ctx, doHTTPRequestParams{
		Method: http.MethodGet,
		Path:   "",
	})
}

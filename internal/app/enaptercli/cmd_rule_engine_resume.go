package enaptercli

import (
	"context"
	"net/http"

	"github.com/urfave/cli/v2"
)

type cmdRuleEngineResume struct {
	cmdRuleEngine
}

func buildCmdRuleEngineResume() *cli.Command {
	cmd := &cmdRuleEngineResume{}
	return &cli.Command{
		Name:               "resume",
		Usage:              "Resume execution of rules",
		CustomHelpTemplate: cmd.HelpTemplate(),
		Flags:              cmd.Flags(),
		Before:             cmd.Before,
		Action: func(cliCtx *cli.Context) error {
			return cmd.do(cliCtx.Context)
		},
	}
}

func (c *cmdRuleEngineResume) do(ctx context.Context) error {
	return c.doHTTPRequest(ctx, doHTTPRequestParams{
		Method: http.MethodPost,
		Path:   "/resume",
	})
}

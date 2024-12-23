package enaptercli

import (
	"context"
	"net/http"

	"github.com/urfave/cli/v3"
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
		Action: func(ctx context.Context, _ *cli.Command) error {
			return cmd.do(ctx)
		},
	}
}

func (c *cmdRuleEngineResume) do(ctx context.Context) error {
	return c.doHTTPRequest(ctx, doHTTPRequestParams{
		Method: http.MethodPost,
		Path:   "/resume",
	})
}

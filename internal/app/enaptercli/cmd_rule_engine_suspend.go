package enaptercli

import (
	"context"
	"net/http"

	"github.com/urfave/cli/v3"
)

type cmdRuleEngineSuspend struct {
	cmdRuleEngine
}

func buildCmdRuleEngineSuspend() *cli.Command {
	cmd := &cmdRuleEngineSuspend{}
	return &cli.Command{
		Name:               "suspend",
		Usage:              "Suspend execution of rules",
		CustomHelpTemplate: cmd.HelpTemplate(),
		Flags:              cmd.Flags(),
		Before:             cmd.Before,
		Action: func(ctx context.Context, _ *cli.Command) error {
			return cmd.do(ctx)
		},
	}
}

func (c *cmdRuleEngineSuspend) do(ctx context.Context) error {
	return c.doHTTPRequest(ctx, doHTTPRequestParams{
		Method: http.MethodPost,
		Path:   "/suspend",
	})
}

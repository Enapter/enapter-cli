package enaptercli

import (
	"context"
	"net/http"

	"github.com/urfave/cli/v2"
)

type cmdRuleEngineSuspend struct {
	cmdRuleEngine
}

func buildCmdRuleEngineSuspend() *cli.Command {
	cmd := &cmdRuleEngineSuspend{}
	return &cli.Command{
		Name:               "suspend",
		Usage:              "Suspend execution of rules",
		CustomHelpTemplate: cmd.CommandHelpTemplate(),
		Flags:              cmd.Flags(),
		Before:             cmd.Before,
		Action: func(cliCtx *cli.Context) error {
			return cmd.do(cliCtx.Context)
		},
	}
}

func (c *cmdRuleEngineSuspend) do(ctx context.Context) error {
	return c.doHTTPRequest(ctx, doHTTPRequestParams{
		Method: http.MethodPost,
		Path:   "/suspend",
	})
}

package enaptercli

import (
	"context"
	"net/http"

	"github.com/urfave/cli/v3"
)

type cmdRuleEngineGet struct {
	cmdRuleEngine
}

func buildCmdRuleEngineGet() *cli.Command {
	cmd := &cmdRuleEngineGet{}
	return &cli.Command{
		Name:               "get",
		Usage:              "Retrieve the rule engine",
		CustomHelpTemplate: cmd.CommandHelpTemplate(),
		Flags:              cmd.Flags(),
		Before:             cmd.Before,
		Action: func(ctx context.Context, _ *cli.Command) error {
			return cmd.do(ctx)
		},
	}
}

func (c *cmdRuleEngineGet) do(ctx context.Context) error {
	return c.doHTTPRequest(ctx, doHTTPRequestParams{
		Method: http.MethodGet,
		Path:   "",
	})
}

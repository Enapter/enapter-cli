package enaptercli

import (
	"context"
	"net/http"

	"github.com/urfave/cli/v2"
)

type cmdDeviceCommandList struct {
	cmdDeviceCommand
}

func buildCmdDeviceCommandList() *cli.Command {
	cmd := &cmdDeviceCommandList{}
	return &cli.Command{
		Name:               "list",
		Usage:              "List device command executions",
		CustomHelpTemplate: cmd.CommandHelpTemplate(),
		Flags:              cmd.Flags(),
		Before:             cmd.Before,
		Action: func(cliCtx *cli.Context) error {
			return cmd.do(cliCtx.Context)
		},
	}
}

func (c *cmdDeviceCommandList) do(ctx context.Context) error {
	return c.doHTTPRequest(ctx, doHTTPRequestParams{
		Method: http.MethodGet,
	})
}

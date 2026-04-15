package enaptercli

import (
	"context"
	"net/http"

	"github.com/urfave/cli/v3"
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
		Action: func(ctx context.Context, _ *cli.Command) error {
			return cmd.do(ctx)
		},
	}
}

func (c *cmdDeviceCommandList) do(ctx context.Context) error {
	return c.doHTTPRequest(ctx, doHTTPRequestParams{
		Method: http.MethodGet,
	})
}

package enaptercli

import (
	"context"
	"net/http"

	"github.com/urfave/cli/v2"
)

type cmdDeviceDelete struct {
	cmdDevice
	deviceID string
}

func buildCmdDeviceDelete() *cli.Command {
	cmd := &cmdDeviceDelete{}
	return &cli.Command{
		Name:               "delete",
		Usage:              "Delete a device",
		CustomHelpTemplate: cmd.CommandHelpTemplate(),
		Flags:              cmd.Flags(),
		Before:             cmd.Before,
		Action: func(cliCtx *cli.Context) error {
			return cmd.do(cliCtx.Context)
		},
	}
}

func (c *cmdDeviceDelete) Flags() []cli.Flag {
	flags := c.cmdDevice.Flags()
	return append(flags,
		&cli.StringFlag{
			Name:        "device-id",
			Aliases:     []string{"d"},
			Usage:       "device ID",
			Destination: &c.deviceID,
			Required:    true,
		},
	)
}

func (c *cmdDeviceDelete) do(ctx context.Context) error {
	return c.doHTTPRequest(ctx, doHTTPRequestParams{
		Method: http.MethodDelete,
		Path:   "/" + c.deviceID,
	})
}

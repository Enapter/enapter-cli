package enaptercli

import (
	"context"
	"net/http"

	"github.com/urfave/cli/v3"
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
		Action: func(ctx context.Context, _ *cli.Command) error {
			return cmd.do(ctx)
		},
	}
}

func (c *cmdDeviceDelete) Flags() []cli.Flag {
	flags := c.cmdDevice.Flags()
	return append(flags,
		&cli.StringFlag{
			Name:        "device-id",
			Aliases:     []string{"d"},
			Usage:       "Device ID",
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

package enaptercli

import (
	"context"
	"net/http"

	"github.com/urfave/cli/v2"
)

type cmdDevicesDelete struct {
	cmdDevices
	deviceID string
}

func buildCmdDevicesDelete() *cli.Command {
	cmd := &cmdDevicesDelete{}
	return &cli.Command{
		Name:               "delete",
		Usage:              "Delete a device",
		CustomHelpTemplate: cmd.HelpTemplate(),
		Flags:              cmd.Flags(),
		Before:             cmd.Before,
		Action: func(cliCtx *cli.Context) error {
			return cmd.do(cliCtx.Context)
		},
	}
}

func (c *cmdDevicesDelete) Flags() []cli.Flag {
	flags := c.cmdDevices.Flags()
	return append(flags,
		&cli.StringFlag{
			Name:        "device-id",
			Usage:       "device ID",
			Destination: &c.deviceID,
			Required:    true,
		},
	)
}

func (c *cmdDevicesDelete) do(ctx context.Context) error {
	return c.doHTTPRequest(ctx, doHTTPRequestParams{
		Method: http.MethodDelete,
		Path:   "/" + c.deviceID,
	})
}

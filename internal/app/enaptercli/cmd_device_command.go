package enaptercli

import (
	"context"
	"fmt"
	"net/url"

	"github.com/urfave/cli/v2"
)

type cmdDeviceCommand struct {
	cmdDevice
	deviceID string
}

func buildCmdDeviceCommand() *cli.Command {
	cmd := &cmdDeviceCommand{}
	return &cli.Command{
		Name:               "command",
		Usage:              "Manage device commands",
		CustomHelpTemplate: cmd.SubcommandHelpTemplate(),
		Subcommands: []*cli.Command{
			buildCmdDeviceCommandExecute(),
			buildCmdDeviceCommandList(),
			buildCmdDeviceCommandGet(),
		},
	}
}

func (c *cmdDeviceCommand) Flags() []cli.Flag {
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

func (c *cmdDeviceCommand) doHTTPRequest(ctx context.Context, p doHTTPRequestParams) error {
	path, err := url.JoinPath(c.deviceID, "command_executions", p.Path)
	if err != nil {
		return fmt.Errorf("join path: %w", err)
	}
	p.Path = path
	return c.cmdDevice.doHTTPRequest(ctx, p)
}

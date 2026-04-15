package enaptercli

import (
	"context"
	"fmt"
	"net/url"

	"github.com/urfave/cli/v3"
)

type cmdDeviceCommunicationConfig struct {
	cmdDevice
	deviceID string
}

func buildCmdDeviceCommunicationConfig() *cli.Command {
	cmd := &cmdDeviceCommunicationConfig{}
	return &cli.Command{
		Name:               "communication-config",
		Usage:              "Manage device communication config",
		CustomHelpTemplate: cmd.SubcommandHelpTemplate(),
		Commands: []*cli.Command{
			buildCmdDeviceCommunicationConfigGenerate(),
		},
	}
}

func (c *cmdDeviceCommunicationConfig) Flags() []cli.Flag {
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

func (c *cmdDeviceCommunicationConfig) doHTTPRequest(ctx context.Context, p doHTTPRequestParams) error {
	path, err := url.JoinPath(c.deviceID, p.Path)
	if err != nil {
		return fmt.Errorf("join path: %w", err)
	}
	p.Path = path
	return c.cmdDevice.doHTTPRequest(ctx, p)
}

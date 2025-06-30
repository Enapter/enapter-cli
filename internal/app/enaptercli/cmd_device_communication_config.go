package enaptercli

import (
	"context"
	"fmt"
	"net/url"

	"github.com/urfave/cli/v2"
)

type cmdDeviceCommunicationConfig struct {
	cmdDevices
	deviceID string
}

func buildCmdDeviceCommunicationConfig() *cli.Command {
	cmd := &cmdDeviceCommunicationConfig{}
	return &cli.Command{
		Name:               "communication-config",
		Usage:              "Manage device communication config",
		CustomHelpTemplate: cmd.SubcommandHelpTemplate(),
		Subcommands: []*cli.Command{
			buildCmdDeviceCommunicationConfigGenerate(),
		},
	}
}

func (c *cmdDeviceCommunicationConfig) Flags() []cli.Flag {
	flags := c.cmdDevices.Flags()
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

func (c *cmdDeviceCommunicationConfig) doHTTPRequest(ctx context.Context, p doHTTPRequestParams) error {
	path, err := url.JoinPath(c.deviceID, p.Path)
	if err != nil {
		return fmt.Errorf("join path: %w", err)
	}
	p.Path = path
	return c.cmdDevices.doHTTPRequest(ctx, p)
}

package enaptercli

import (
	"context"
	"fmt"
	"net/url"

	"github.com/urfave/cli/v2"
)

type cmdDeviceExecution struct {
	cmdDevices
	deviceID string
}

func buildCmdDeviceExecution() *cli.Command {
	return &cli.Command{
		Name:  "execution",
		Usage: "Manage device command executions",
		Subcommands: []*cli.Command{
			buildCmdDeviceExecutionList(),
			buildCmdDeviceExecutionInspect(),
		},
	}
}

func (c *cmdDeviceExecution) Flags() []cli.Flag {
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

func (c *cmdDeviceExecution) doHTTPRequest(ctx context.Context, p doHTTPRequestParams) error {
	path, err := url.JoinPath(c.deviceID, "command_executions", p.Path)
	if err != nil {
		return fmt.Errorf("join path: %w", err)
	}
	p.Path = path
	return c.cmdDevices.doHTTPRequest(ctx, p)
}

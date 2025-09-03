package enaptercli

import (
	"context"
	"fmt"
	"net/url"

	"github.com/urfave/cli/v2"
)

type cmdDevices struct {
	cmdBase
}

func buildCmdDevices() *cli.Command {
	cmd := &cmdDevices{}
	return &cli.Command{
		Name:               "device",
		Usage:              "Manage devices",
		CustomHelpTemplate: cmd.SubcommandHelpTemplate(),
		Subcommands: []*cli.Command{
			buildCmdDeviceCreate(),
			buildCmdDevicesList(),
			buildCmdDevicesGet(),
			buildCmdDevicesChangeBlueprint(),
			buildCmdDevicesLogs(),
			buildCmdDevicesDelete(),
			buildCmdDeviceCommand(),
			buildCmdDeviceTelemetry(),
			buildCmdDeviceCommunicationConfig(),
		},
	}
}

func (c *cmdDevices) doHTTPRequest(ctx context.Context, p doHTTPRequestParams) error {
	path, err := url.JoinPath("/devices", p.Path)
	if err != nil {
		return fmt.Errorf("join path: %w", err)
	}
	p.Path = path
	return c.cmdBase.doHTTPRequest(ctx, p)
}

func (c *cmdDevices) validateExpandFlag(cliCtx *cli.Context) error {
	return validateExpandFlag(cliCtx, c.supportedExpandFields())
}

func (c *cmdDevices) supportedExpandFields() []string {
	return []string{"connectivity", "manifest", "properties", "communication", "site"}
}

package enaptercli

import (
	"context"
	"fmt"
	"net/url"

	"github.com/urfave/cli/v2"
)

type cmdDevice struct {
	cmdBase
}

func buildCmdDevice() *cli.Command {
	cmd := &cmdDevice{}
	return &cli.Command{
		Name:               "device",
		Usage:              "Manage devices",
		CustomHelpTemplate: cmd.SubcommandHelpTemplate(),
		Subcommands: []*cli.Command{
			buildCmdDeviceCreate(),
			buildCmdDeviceList(),
			buildCmdDeviceGet(),
			buildCmdDeviceChangeBlueprint(),
			buildCmdDeviceLogs(),
			buildCmdDeviceUpdate(),
			buildCmdDeviceDelete(),
			buildCmdDeviceCommand(),
			buildCmdDeviceTelemetry(),
			buildCmdDeviceCommunicationConfig(),
		},
	}
}

func (c *cmdDevice) doHTTPRequest(ctx context.Context, p doHTTPRequestParams) error {
	path, err := url.JoinPath("/devices", p.Path)
	if err != nil {
		return fmt.Errorf("join path: %w", err)
	}
	p.Path = path
	return c.cmdBase.doHTTPRequest(ctx, p)
}

func (c *cmdDevice) validateExpandFlag(cliCtx *cli.Context) error {
	return validateExpandFlag(cliCtx, c.supportedExpandFields())
}

func (c *cmdDevice) supportedExpandFields() []string {
	return []string{"connectivity", "manifest", "properties", "communication", "site"}
}

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
			buildCmdDevicesList(),
			buildCmdDevicesGet(),
			buildCmdDevicesAssignBlueprint(),
			buildCmdDevicesLogs(),
			buildCmdDevicesDelete(),
			buildCmdDeviceCommand(),
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
	supportedFields := []string{"connectivity", "manifest", "properties", "communication_info", "site"}
	return validateExpandFlag(cliCtx, supportedFields)
}

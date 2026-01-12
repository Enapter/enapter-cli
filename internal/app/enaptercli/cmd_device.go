package enaptercli

import (
	"context"
	"errors"
	"fmt"
	"net/url"

	"github.com/urfave/cli/v2"
)

type cmdDevice struct {
	cmdBase
	siteID string
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
			buildCmdDeviceStream(),
			buildCmdDeviceCommunicationConfig(),
			buildCmdDeviceRunTerminal(),
		},
	}
}

func (c *cmdDevice) Flags() []cli.Flag {
	flags := c.cmdBase.Flags()
	return append(flags, &cli.StringFlag{
		Name:        "site-id",
		Usage:       "Site ID",
		Destination: &c.siteID,
	})
}

func (c *cmdDevice) doHTTPRequest(ctx context.Context, p doHTTPRequestParams) error {
	path, err := c.buildPath(p.Path)
	if err != nil {
		return err
	}
	p.Path = path
	return c.cmdBase.doHTTPRequest(ctx, p)
}

func (c *cmdDevice) runWebSocket(ctx context.Context, p runWebSocketParams) error {
	path, err := c.buildPath(p.Path)
	if err != nil {
		return err
	}
	p.Path = path
	return c.cmdBase.runWebSocket(ctx, p)
}

func (c *cmdDevice) validateExpandFlag(cliCtx *cli.Context) error {
	return validateExpandFlag(cliCtx, c.supportedExpandFields())
}

func (c *cmdDevice) supportedExpandFields() []string {
	return []string{"connectivity", "manifest", "properties", "communication", "site"}
}

func (c *cmdDevice) buildPath(p string) (string, error) {
	path, err := url.JoinPath("/devices", p)
	if err != nil {
		return "", fmt.Errorf("join path: %w", err)
	}

	siteID, err := c.chooseSiteID(c.siteID)
	if err != nil {
		if errors.Is(err, errSiteIDMissing) {
			return path, nil
		}
		return "", err
	}

	path, err = url.JoinPath("/sites", siteID, path)
	if err != nil {
		return "", fmt.Errorf("join path: %w", err)
	}

	return path, nil
}

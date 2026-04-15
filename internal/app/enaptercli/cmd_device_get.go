package enaptercli

import (
	"context"
	"net/http"
	"net/url"
	"strings"

	"github.com/urfave/cli/v3"
)

type cmdDeviceGet struct {
	cmdDevice
	deviceID string
	expand   []string
}

func buildCmdDeviceGet() *cli.Command {
	cmd := &cmdDeviceGet{}
	return &cli.Command{
		Name:               "get",
		Usage:              "Retrieve device information",
		CustomHelpTemplate: cmd.CommandHelpTemplate(),
		Flags:              cmd.Flags(),
		Before:             cmd.Before,
		Action: func(ctx context.Context, _ *cli.Command) error {
			return cmd.do(ctx)
		},
	}
}

func (c *cmdDeviceGet) Flags() []cli.Flag {
	flags := c.cmdDevice.Flags()
	return append(flags, &cli.StringFlag{
		Name:        "device-id",
		Aliases:     []string{"d"},
		Usage:       "Device ID",
		Destination: &c.deviceID,
		Required:    true,
	}, &cli.StringSliceFlag{
		Name: "expand",
		Usage: "Comma-separated list of expanded device information (supported values: " +
			strings.Join(c.supportedExpandFields(), ", ") + ")",
		Destination: &c.expand,
	})
}

func (c *cmdDeviceGet) Before(ctx context.Context, cmd *cli.Command) (context.Context, error) {
	ctx, err := c.cmdDevice.Before(ctx, cmd)
	if err != nil {
		return ctx, err
	}
	return ctx, c.validateExpandFlag(cmd)
}

func (c *cmdDeviceGet) do(ctx context.Context) error {
	query := url.Values{}
	if len(c.expand) != 0 {
		query.Set("expand", strings.Join(c.expand, ","))
	}
	return c.doHTTPRequest(ctx, doHTTPRequestParams{
		Method: http.MethodGet,
		Path:   "/" + c.deviceID,
		Query:  query,
	})
}

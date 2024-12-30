package enaptercli

import (
	"context"
	"net/http"
	"net/url"
	"strings"

	"github.com/urfave/cli/v2"
)

type cmdDevicesInspect struct {
	cmdDevices
	deviceID string
	expand   []string
}

func buildCmdDevicesInspect() *cli.Command {
	cmd := &cmdDevicesInspect{}
	return &cli.Command{
		Name:               "inspect",
		Usage:              "Inspect a devices",
		CustomHelpTemplate: cmd.HelpTemplate(),
		Flags:              cmd.Flags(),
		Before:             cmd.Before,
		Action: func(cliCtx *cli.Context) error {
			return cmd.do(cliCtx.Context)
		},
	}
}

func (c *cmdDevicesInspect) Flags() []cli.Flag {
	flags := c.cmdDevices.Flags()
	return append(flags, &cli.StringFlag{
		Name:        "device-id",
		Aliases:     []string{"d"},
		Usage:       "device ID",
		Destination: &c.deviceID,
		Required:    true,
	}, &cli.MultiStringFlag{
		Target: &cli.StringSliceFlag{
			Name:  "expand",
			Usage: "coma separated list of expanded device info",
		},
		Destination: &c.expand,
	})
}

func (c *cmdDevicesInspect) Before(cliCtx *cli.Context) error {
	if err := c.cmdDevices.Before(cliCtx); err != nil {
		return err
	}
	return c.validateExpandFlag(cliCtx)
}

func (c *cmdDevicesInspect) do(ctx context.Context) error {
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

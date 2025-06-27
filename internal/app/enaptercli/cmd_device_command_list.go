package enaptercli

import (
	"context"
	"net/http"
	"net/url"
	"strings"

	"github.com/urfave/cli/v2"
)

type cmdDeviceCommandList struct {
	cmdDeviceCommand
	expand []string
}

func buildCmdDeviceCommandList() *cli.Command {
	cmd := &cmdDeviceCommandList{}
	return &cli.Command{
		Name:               "list",
		Usage:              "List device command executions",
		CustomHelpTemplate: cmd.CommandHelpTemplate(),
		Flags:              cmd.Flags(),
		Before:             cmd.Before,
		Action: func(cliCtx *cli.Context) error {
			return cmd.do(cliCtx.Context)
		},
	}
}

func (c *cmdDeviceCommandList) Flags() []cli.Flag {
	flags := c.cmdDeviceCommand.Flags()
	return append(flags,
		&cli.MultiStringFlag{
			Target: &cli.StringSliceFlag{
				Name:  "expand",
				Usage: "coma separated list of expanded options",
			},
			Destination: &c.expand,
		},
	)
}

func (c *cmdDeviceCommandList) Before(cliCtx *cli.Context) error {
	if err := c.cmdDevices.Before(cliCtx); err != nil {
		return err
	}
	return validateExpandFlag(cliCtx, []string{"ephemeral"})
}

func (c *cmdDeviceCommandList) do(ctx context.Context) error {
	query := url.Values{}
	if len(c.expand) != 0 {
		query.Set("expand", strings.Join(c.expand, ","))
	}
	return c.doHTTPRequest(ctx, doHTTPRequestParams{
		Method: http.MethodGet,
		Query:  query,
	})
}

package enaptercli

import (
	"context"
	"net/http"
	"net/url"
	"strings"

	"github.com/urfave/cli/v2"
)

type cmdDeviceExecutionList struct {
	cmdDeviceExecution
	expand []string
}

func buildCmdDeviceExecutionList() *cli.Command {
	cmd := &cmdDeviceExecutionList{}
	return &cli.Command{
		Name:               "list",
		Usage:              "List device command executions",
		CustomHelpTemplate: cmd.HelpTemplate(),
		Flags:              cmd.Flags(),
		Before:             cmd.Before,
		Action: func(cliCtx *cli.Context) error {
			return cmd.do(cliCtx.Context)
		},
	}
}

func (c *cmdDeviceExecutionList) Flags() []cli.Flag {
	flags := c.cmdDeviceExecution.Flags()
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

func (c *cmdDeviceExecutionList) Before(cliCtx *cli.Context) error {
	if err := c.cmdDevices.Before(cliCtx); err != nil {
		return err
	}
	return validateExpandFlag(cliCtx, []string{"ephemeral"})
}

func (c *cmdDeviceExecutionList) do(ctx context.Context) error {
	query := url.Values{}
	if len(c.expand) != 0 {
		query.Set("expand", strings.Join(c.expand, ","))
	}
	return c.doHTTPRequest(ctx, doHTTPRequestParams{
		Method: http.MethodGet,
		Query:  query,
	})
}

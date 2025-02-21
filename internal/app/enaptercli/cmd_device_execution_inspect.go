package enaptercli

import (
	"context"
	"net/http"
	"net/url"
	"strings"

	"github.com/urfave/cli/v2"
)

type cmdDeviceExecutionInspect struct {
	cmdDeviceExecution
	executionID string
	expand      []string
}

func buildCmdDeviceExecutionInspect() *cli.Command {
	cmd := &cmdDeviceExecutionInspect{}
	return &cli.Command{
		Name:               "inspect",
		Usage:              "Inspect a device command execution",
		CustomHelpTemplate: cmd.HelpTemplate(),
		Flags:              cmd.Flags(),
		Before:             cmd.Before,
		Action: func(cliCtx *cli.Context) error {
			return cmd.do(cliCtx.Context)
		},
	}
}

func (c *cmdDeviceExecutionInspect) Flags() []cli.Flag {
	flags := c.cmdDeviceExecution.Flags()
	return append(flags,
		&cli.StringFlag{
			Name:        "execution-id",
			Usage:       "execution ID",
			Destination: &c.executionID,
			Required:    true,
		}, &cli.MultiStringFlag{
			Target: &cli.StringSliceFlag{
				Name:  "expand",
				Usage: "coma separated list of expanded options",
			},
			Destination: &c.expand,
		},
	)
}

func (c *cmdDeviceExecutionInspect) Before(cliCtx *cli.Context) error {
	if err := c.cmdDevices.Before(cliCtx); err != nil {
		return err
	}
	return validateExpandFlag(cliCtx, []string{"log"})
}

func (c *cmdDeviceExecutionInspect) do(ctx context.Context) error {
	query := url.Values{}
	if len(c.expand) != 0 {
		query.Set("expand", strings.Join(c.expand, ","))
	}
	return c.doHTTPRequest(ctx, doHTTPRequestParams{
		Method: http.MethodGet,
		Path:   "/" + c.executionID,
		Query:  query,
	})
}

package enaptercli

import (
	"context"
	"net/http"
	"net/url"
	"strings"

	"github.com/urfave/cli/v2"
)

type cmdDeviceCommandGet struct {
	cmdDeviceCommand
	executionID string
	expand      []string
}

func buildCmdDeviceCommandGet() *cli.Command {
	cmd := &cmdDeviceCommandGet{}
	return &cli.Command{
		Name:               "get",
		Usage:              "Get a device command execution",
		CustomHelpTemplate: cmd.CommandHelpTemplate(),
		Flags:              cmd.Flags(),
		Before:             cmd.Before,
		Action: func(cliCtx *cli.Context) error {
			return cmd.do(cliCtx.Context)
		},
	}
}

func (c *cmdDeviceCommandGet) Flags() []cli.Flag {
	flags := c.cmdDeviceCommand.Flags()
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

func (c *cmdDeviceCommandGet) Before(cliCtx *cli.Context) error {
	if err := c.cmdDevices.Before(cliCtx); err != nil {
		return err
	}
	return validateExpandFlag(cliCtx, []string{"log"})
}

func (c *cmdDeviceCommandGet) do(ctx context.Context) error {
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

package enaptercli

import (
	"context"
	"net/http"
	"net/url"
	"strings"

	"github.com/urfave/cli/v3"
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
		Usage:              "Retrieve a device command execution",
		CustomHelpTemplate: cmd.CommandHelpTemplate(),
		Flags:              cmd.Flags(),
		Before:             cmd.Before,
		Action: func(ctx context.Context, _ *cli.Command) error {
			return cmd.do(ctx)
		},
	}
}

func (c *cmdDeviceCommandGet) Flags() []cli.Flag {
	flags := c.cmdDeviceCommand.Flags()
	return append(flags,
		&cli.StringFlag{
			Name:        "execution-id",
			Usage:       "Execution ID",
			Destination: &c.executionID,
			Required:    true,
		}, &cli.StringSliceFlag{
			Name:        "expand",
			Usage:       "Comma-separated list of expanded options (supported values: log)",
			Destination: &c.expand,
		},
	)
}

func (c *cmdDeviceCommandGet) Before(ctx context.Context, cmd *cli.Command) (context.Context, error) {
	ctx, err := c.cmdDevice.Before(ctx, cmd)
	if err != nil {
		return ctx, err
	}
	return ctx, validateExpandFlag(cmd, []string{"log"})
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

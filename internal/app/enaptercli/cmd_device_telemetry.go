package enaptercli

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/urfave/cli/v2"
)

type cmdDeviceTelemetry struct {
	cmdDevices
	deviceID string
	follow   bool
}

func buildCmdDeviceTelemetry() *cli.Command {
	cmd := &cmdDeviceTelemetry{}
	return &cli.Command{
		Name:               "telemetry",
		Usage:              "Show device telemetry",
		CustomHelpTemplate: cmd.CommandHelpTemplate(),
		Flags:              cmd.Flags(),
		Before:             cmd.Before,
		Action: func(cliCtx *cli.Context) error {
			return cmd.do(cliCtx.Context)
		},
	}
}

func (c *cmdDeviceTelemetry) Flags() []cli.Flag {
	flags := c.cmdDevices.Flags()
	return append(flags, &cli.StringFlag{
		Name:        "device-id",
		Aliases:     []string{"d"},
		Usage:       "device ID",
		Destination: &c.deviceID,
		Required:    true,
	}, &cli.BoolFlag{
		Name:        "follow",
		Aliases:     []string{"f"},
		Usage:       "follow the telemetry output",
		Destination: &c.follow,
	})
}

func (c *cmdDeviceTelemetry) do(ctx context.Context) error {
	if c.follow {
		return c.doFollow(ctx)
	}
	return c.doList(ctx)
}

func (c *cmdDeviceTelemetry) doFollow(ctx context.Context) error {
	return c.runWebSocket(ctx, runWebSocketParams{
		Path: "/devices/" + c.deviceID + "/telemetry",
		RespProcessor: func(r io.Reader) error {
			payload, err := io.ReadAll(r)
			if err != nil {
				return fmt.Errorf("read response: %w", err)
			}
			fmt.Fprintln(c.writer, strings.TrimSpace(string(payload)))
			return nil
		},
	})
}

func (c *cmdDeviceTelemetry) doList(ctx context.Context) error {
	return c.doHTTPRequest(ctx, doHTTPRequestParams{
		Method: http.MethodGet,
		Path:   "/" + c.deviceID + "/telemetry",
		//nolint:bodyclose //body is closed in doHTTPRequest
		RespProcessor: okRespBodyProcessor(func(r io.Reader) error {
			payload, err := io.ReadAll(r)
			if err != nil {
				return fmt.Errorf("read response: %w", err)
			}
			fmt.Fprintln(c.writer, strings.TrimSpace(string(payload)))
			return nil
		}),
	})
}
